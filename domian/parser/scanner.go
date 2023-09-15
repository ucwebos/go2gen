package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/xbitgo/core/tools/tool_str"
)

const (
	ParseTypeWatch   = 1
	ParseTypeImpl    = 2
	ParseTypeDo      = 3
	ParseTypeHandler = 4
)

func Scan(pwd string, parseType int) (ips *IParser, err error) {
	ips = &IParser{
		Pwd:          pwd,
		Package:      tool_str.LastPwdStr(pwd),
		INFList:      make(map[string]INF, 0),
		StructList:   make(map[string]XST, 0),
		OtherStruct:  make(map[string]XST, 0),
		ConstStrList: make(map[string]string),
		BindFuncMap:  map[string]map[string]XMethod{},
		FuncList:     map[string]XMethod{},
		NewFuncList:  map[string]XMethod{},
		EntryModules: []EntryModule{},
		ParseType:    parseType,
	}
	// 遍历文件夹解析go代码
	fileInfos, err := os.ReadDir(pwd)
	if err != nil {
		return
	}
	for _, fi := range fileInfos {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
			iPwd := pwd + "/" + fi.Name()
			//fmt.Println(iPwd)
			err := ips.ParseFile(iPwd)
			if err != nil {
				log.Printf("ips.ParseFile [%s] err: %v \n", iPwd, err)
				continue
			}
		}
	}
	ips.LoadBinds()
	return
}

func (ips *IParser) LoadBinds() {
	// 绑定方法入结构体
	for sName, sts := range ips.StructList {
		smths := make(map[string]XMethod, 0)
		point := false
		if mths, ok := ips.BindFuncMap[sName]; ok {
			smths = mths
		}
		if mths, ok := ips.BindFuncMap["*"+sName]; ok {
			point = true
			smths = mths
		}
		for _, cs := range sts.CST {
			if cMths, ok := ips.BindFuncMap[cs]; ok {
				for s, method := range cMths {
					smths[s] = method
				}
			}
			if cMths, ok := ips.BindFuncMap["*"+cs]; ok {
				point = true
				for s, method := range cMths {
					smths[s] = method
				}
			}
		}
		if len(smths) == 0 {
			continue
		}
		sts.MPoint = point
		sts.Methods = smths
		for _, method := range smths {
			sts.ShortName = method.ImplName
			break
		}
		ips.StructList[sName] = sts
	}
}

var reImpl, _ = regexp.Compile(`@IMPL\[([\w|.]+)]`)
var reDi, _ = regexp.Compile(`@DI\[([\w|.]+)]`)

var reModule, _ = regexp.Compile(`(\S+)\(([\w|.]+)\)`)
var reMiddleware, _ = regexp.Compile(`@MIDDLEWARE\[([\w|.]+)]`)

func (ips *IParser) ParseFile(pwd string) error {
	r, err := os.Open(pwd)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, pwd, r, parser.ParseComments)
	if err != nil {
		return err
	}
	// ----------------DEBUG------------------
	//if strings.Contains(pwd, "_impl.go") {
	//	txt, err := os.OpenFile("ast.txt", os.O_WRONLY|os.O_CREATE, 0666)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	err = ast.Fprint(txt, fset, f, ast.NotNilFilter)
	//	fmt.Println(err)
	//}
	// ----------------DEBUG------------------

	imports := make([]string, 0)
	for _, decl := range f.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			switch x.Tok {
			case token.IMPORT:
				for _, spec := range x.Specs {
					xSpec := spec.(*ast.ImportSpec)
					imports = append(imports, xSpec.Path.Value)
				}
			case token.TYPE:
				if ips.ParseType == ParseTypeHandler {
					module := EntryModule{
						Name:       "",
						Key:        "",
						Middleware: []string{},
						WithCommon: false,
						FuncList:   make([]*EntryModuleFunc, 0),
					}
					if x.Doc != nil && len(x.Doc.List) >= 0 {
						for _, comment := range x.Doc.List {
							if strings.Index(comment.Text, "@COMMON") > 0 {
								module.WithCommon = true
							}
							r := reModule.FindStringSubmatch(comment.Text)
							if len(r) == 3 {
								module.Name = r[1]
								module.Key = r[2]
							}
							if strings.Contains(comment.Text, "@MIDDLEWARE[") {
								r := reMiddleware.FindStringSubmatch(comment.Text)
								if len(r) == 2 {
									list := strings.Split(r[1], ",")
									module.Middleware = list
								}
							}
						}
					}
					tmpFuncs := map[string]*EntryModuleFunc{}
					for i, spec := range x.Specs {
						var isReq = false
						xSpec := spec.(*ast.TypeSpec)
						nameDoc := xSpec.Doc.Text()
						name := xSpec.Name.Name
						key := ""
						if !strings.HasSuffix(name, "Req") && !strings.HasSuffix(name, "Resp") {
							continue
						}
						if strings.HasSuffix(name, "Req") {
							isReq = true
							key = strings.TrimSuffix(name, "Req")
						} else {
							key = strings.TrimSuffix(name, "Resp")
						}

						switch xt := xSpec.Type.(type) {
						case *ast.StructType:
							xst := XST{
								Imports:   imports,
								File:      pwd,
								Name:      name,
								CST:       make([]string, 0),
								Methods:   make(map[string]XMethod, 0),
								FieldList: make(map[string]XField, 0),
							}
							fields, child := getStructField(xt)
							xst.CST = child
							xst.FieldList = fields
							_, ok := tmpFuncs[key]
							if !ok {
								tmpFuncs[key] = &EntryModuleFunc{
									idx:   i,
									Name:  strings.TrimSuffix(nameDoc, "\n"),
									Key:   key,
									KeyLi: strings.ReplaceAll(tool_str.ToSnakeCase(key), "_", "-"),
								}
							}
							if isReq {
								tmpFuncs[key].Request = &EntryModuleFuncReq{
									Name: name,
									XST:  xst,
								}
							} else {
								tmpFuncs[key].Response = &EntryModuleFuncResp{
									Name: name,
									XST:  xst,
								}
							}
						}
					}

					for _, moduleFunc := range tmpFuncs {
						module.FuncList = append(module.FuncList, moduleFunc)
					}
					ips.EntryModules = append(ips.EntryModules, module)
				} else {
					var (
						used       = true
						ImplINF    = ""
						GIName     = ""
						GI         = false
						noDeleteAT = false
					)
					if x.Doc != nil && len(x.Doc.List) >= 0 {
						for _, comment := range x.Doc.List {
							if strings.Index(comment.Text, "@IGNORE") > 0 {
								used = false
							}
							if strings.Contains(comment.Text, "@IMPL[") {
								r := reImpl.FindStringSubmatch(comment.Text)
								if len(r) > 1 {
									ImplINF = r[1]
								}
							}
							if strings.Index(comment.Text, "@NODELETEAT") > 0 {
								noDeleteAT = true
							}
							if strings.Contains(comment.Text, "@GI") {
								GI = true
								r := reDi.FindStringSubmatch(comment.Text)
								if len(r) > 1 {
									GIName = r[1]
								}
							}
						}
					}

					for _, spec := range x.Specs {
						xSpec := spec.(*ast.TypeSpec)
						name := xSpec.Name.Name
						switch xt := xSpec.Type.(type) {
						case *ast.InterfaceType:
							if ips.ParseType != ParseTypeImpl {
								continue
							}
							if !used {
								continue
							}
							inf := INF{
								Imports: imports,
								File:    pwd,
								Name:    name,
								Methods: make(map[string]XMethod, 0),
							}
							inf.Methods = getInterfaceFunc(xt)
							ips.INFList[name] = inf
						case *ast.StructType:
							xst := XST{
								GIName:     GIName,
								GI:         GI,
								NoDeleteAT: noDeleteAT,
								ImplINF:    ImplINF,
								Imports:    imports,
								File:       pwd,
								Name:       name,
								CST:        make([]string, 0),
								Methods:    make(map[string]XMethod, 0),
								FieldList:  make(map[string]XField, 0),
							}
							fields, child := getStructField(xt)
							xst.CST = child
							xst.FieldList = fields
							if !used {
								ips.OtherStruct[name] = xst
								continue
							}
							ips.StructList[name] = xst
						}
					}
				}

			case token.CONST:
				if ips.ParseType != ParseTypeDo {
					continue
				}
				for _, spec := range x.Specs {
					xSpec := spec.(*ast.ValueSpec)
					if len(xSpec.Names) == 1 {
						xName := xSpec.Names[0]
						if strings.HasPrefix(xName.Name, "TableName") {
							if len(xSpec.Values) == 1 {
								ips.ConstStrList[xName.Name] = strings.Trim(xSpec.Values[0].(*ast.BasicLit).Value, `"`)
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if x.Recv != nil {
				var bindName string
				var implName = "impl"
				for _, field := range x.Recv.List {
					bindName, _ = getTypeStr(field.Type)
					if field.Names != nil && len(field.Names) > 0 {
						implName = field.Names[0].Name
					}
				}
				mtd := XMethod{
					ImplName: implName,
					Name:     x.Name.Name,
				}
				mtd.Params, mtd.Results = getFuncArgs(x.Type)
				if _, ok := ips.BindFuncMap[bindName]; !ok {
					ips.BindFuncMap[bindName] = make(map[string]XMethod, 0)
				}
				ips.BindFuncMap[bindName][mtd.Name] = mtd
			} else {
				mtd := XMethod{
					Name: x.Name.Name,
				}
				mtd.Params, mtd.Results = getFuncArgs(x.Type)
				ips.FuncList[x.Name.Name] = mtd
				if strings.HasPrefix(x.Name.Name, "New") {
					mtd := XMethod{
						Name: x.Name.Name,
					}
					mtd.Params, mtd.Results = getFuncArgs(x.Type)
					theTypeName := strings.TrimPrefix(x.Name.Name, "New")
					ips.NewFuncList[theTypeName] = mtd
				}
			}
		}

	}
	return nil
}

func (ips *IParser) ReloadFileForStruct(pwd string) (changeList map[string]XST, err error) {
	changeList = map[string]XST{}

	oStList := make(map[string]XST, 0)
	for s, xst := range ips.StructList {
		oStList[s] = xst
	}
	err = ips.ParseFile(pwd)
	if err != nil {
		return
	}
	nStList := ips.StructList
	for name, xst := range nStList {
		oXst, ok := oStList[name]
		if !ok {
			changeList[name] = xst
			continue
		}
		if xst.IsChangedDI(oXst) {
			changeList[name] = xst
			continue
		}
		if xst.IsChangedField(oXst) {
			changeList[name] = xst
		}
	}
	return
}

func (ips *IParser) ReloadFileForINF(pwd string) (changes []INF, err error) {
	changes = []INF{}

	oINFList := make(map[string]INF, 0)
	for s, inf := range ips.INFList {
		oINFList[s] = inf
	}
	err = ips.ParseFile(pwd)
	if err != nil {
		return
	}
	nINFList := ips.INFList
	for name, inf := range nINFList {
		oInf, ok := oINFList[name]
		if !ok {
			changes = append(changes, inf)
			continue
		}
		if !inf.Equal(oInf) {
			changes = append(changes, inf)
		}
	}
	return
}

func getMapTypeStr(arg *ast.MapType) string {
	var (
		key   string
		value string
	)
	switch x := arg.Key.(type) {
	case *ast.Ident:
		key = x.Name
	}

	value, _ = getTypeStr(arg.Value)

	return fmt.Sprintf("map[%s]%s", key, value)

}

func getSliceTypeStr(arg *ast.ArrayType) string {
	value, _ := getTypeStr(arg.Elt)
	return fmt.Sprintf("[]%s", value)
}

func getTypeStr(arg ast.Node) (string, int) {
	switch _type := arg.(type) {
	case *ast.SelectorExpr:
		tStr := fmt.Sprintf("%s.%s", _type.X, _type.Sel)
		return tStr, STypeStruct
	case *ast.StarExpr:
		switch _type2 := _type.X.(type) {
		case *ast.SelectorExpr:
			return fmt.Sprintf("*%s.%s", _type2.X, _type2.Sel), STypeStruct
		case *ast.Ident:
			if unicode.IsUpper([]rune(_type2.Name)[0]) {
				return fmt.Sprintf("*%s", _type2.Name), STypeStruct
			}
			return fmt.Sprintf("*%s", _type2.Name), STypeBasic
		}
	case *ast.Ident:
		if unicode.IsUpper([]rune(_type.Name)[0]) {
			return _type.Name, STypeStruct
		}
		return _type.Name, STypeBasic
	case *ast.MapType:
		return getMapTypeStr(_type), STypeMap
	case *ast.ArrayType:
		return getSliceTypeStr(_type), STypeSlice
	case *ast.InterfaceType:
		return "interface{}", STypeBasic
	}

	return "", STypeBasic
}

func getFuncArgs(xType *ast.FuncType) (params []XArg, results []XArg) {
	params = make([]XArg, 0)
	results = make([]XArg, 0)
	if xType.Params != nil {
		for i, param := range xType.Params.List {
			arg := XArg{}
			arg.Type, _ = getTypeStr(param.Type)
			if len(param.Names) > 0 {
				arg.Name = param.Names[0].Name
			} else {
				if strings.Contains(arg.Type, "Context") {
					arg.Name = "ctx"
				} else if strings.Contains(arg.Type, "Request") {
					arg.Name = "req"
				} else {
					arg.Name = fmt.Sprintf("arg%d", i)
				}
			}
			params = append(params, arg)
		}
	}
	if xType.Results != nil {
		for _, ret := range xType.Results.List {
			arg := XArg{}
			arg.Type, _ = getTypeStr(ret.Type)
			if len(ret.Names) > 0 {
				arg.Name = ret.Names[0].Name
			}
			results = append(results, arg)
		}
	}
	return params, results
}

func getInterfaceFunc(xType *ast.InterfaceType) map[string]XMethod {
	meths := map[string]XMethod{}
	for idx, mt := range xType.Methods.List {
		comment := ""
		if mt.Doc != nil {
			for _, c := range mt.Doc.List {
				comment = c.Text
			}
		}
		switch xType := mt.Type.(type) {
		case *ast.FuncType:
			xMth := XMethod{
				ImplName: "impl",
				Name:     mt.Names[0].Name,
				Comment:  comment,
				Sort:     idx,
			}
			xMth.Params, xMth.Results = getFuncArgs(xType)
			meths[xMth.Name] = xMth
		case *ast.Ident:
			cInf := xType.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType)
			meths2 := getInterfaceFunc(cInf)
			for s, method := range meths2 {
				if _, ok := meths[s]; ok {
					continue
				}
				meths[s] = method
			}
		}
	}
	return meths
}

func getStructField(xType *ast.StructType) (fields map[string]XField, child []string) {
	fields = make(map[string]XField)
	child = make([]string, 0)
	for idx, fe := range xType.Fields.List {
		if fe.Type == nil {
			continue
		}
		if (fe.Names != nil && len(fe.Names) > 0) || fe.Type != nil {
			name := ""
			if len(fe.Names) > 0 {
				name = fe.Names[0].Name
				if unicode.IsLower([]rune(name)[0]) {
					continue
				}
			}
			fType, sType := getTypeStr(fe.Type)
			xf := XField{
				Name:  name,
				Type:  fType,
				SType: sType,
				Idx:   idx,
			}
			if fe.Tag != nil {
				xf.Tag = fe.Tag.Value
			}
			if fe.Comment != nil {
				xf.Comment = strings.Trim(fe.Comment.Text(), "\n")
			}
			fields[name] = xf
		} else {
			if xt, ok := fe.Type.(*ast.Ident); ok {
				if xt.Obj == nil {
					continue
				}
				if xde, ok := xt.Obj.Decl.(*ast.TypeSpec); ok {
					name := xde.Name.Name
					if xdt, ok := xde.Type.(*ast.StructType); ok {
						cFields, cChild := getStructField(xdt)
						for s, fld := range cFields {
							if _, ok := fields[s]; ok {
								continue
							}
							fields[s] = fld
						}
						child = append(child, cChild...)
					}
					child = append(child, name)
				}
			}
		}
	}
	return fields, child
}
