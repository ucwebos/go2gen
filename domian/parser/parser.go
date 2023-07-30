package parser

import (
	"regexp"
	"strings"
)

const (
	STypeBasic  = 0
	STypeStruct = 1
	STypeSlice  = 2
	STypeMap    = 3
	STypeTime   = 4
)

type IParser struct {
	Pwd          string
	Package      string
	INFList      map[string]INF
	StructList   map[string]XST
	OtherStruct  map[string]XST
	ConstStrList map[string]string
	BindFuncMap  map[string]map[string]XMethod
	FuncList     map[string]XMethod
	NewFuncList  map[string]XMethod
	EntryModules []EntryModule
	ParseType    int
}

type INF struct {
	//AF      *ast.File
	Imports []string
	File    string
	Name    string
	Methods map[string]XMethod
}

func (x INF) Equal(x2 INF) bool {
	for s, method := range x.Methods {
		m2, ok := x2.Methods[s]
		if !ok {
			return false
		}
		if !m2.Equal(method) {
			return false
		}
	}
	return true
}

type XST struct {
	GIName     string // 自定义DI名
	GI         bool   // 是否注册DI
	NoDeleteAT bool
	ImplINF    string // 标注实现接口
	Imports    []string
	File       string             // 所在文件
	Name       string             // 结构体名称
	ShortName  string             // 结构体定义方法时的引用名称
	MPoint     bool               // 是否是使用指针定义方法
	CST        []string           // 子结构体
	Methods    map[string]XMethod // 方法列表
	FieldList  map[string]XField  // 字段列表
}

type XSTChanges struct {
	XST    XST
	Add    []string
	Remove []string
	Modify []string
}

func (x XST) IsChangedDI(ox XST) (changed bool) {
	if ox.GI != x.GI {
		return true
	}
	if ox.GIName != x.GIName {
		return true
	}
	return false
}

func (x XST) IsChangedField(ox XST) (changed bool) {
	for s := range ox.FieldList {
		if _, ok := x.FieldList[s]; !ok {
			return true
		}
	}

	for s, field := range x.FieldList {
		oField, ok := ox.FieldList[s]
		if !ok {
			return true
		}
		if !field.Equal(oField) {
			return true
		}
	}
	return false
}

func (x XST) IsImpl(inf INF) (equal bool, noImplFunc []XMethod, changeFunc []XMethod) {
	equal = true
	noImplFunc = make([]XMethod, 0)
	changeFunc = make([]XMethod, 0)
	for s, method := range inf.Methods {
		sMth, ok := x.Methods[s]
		if !ok {
			equal = false
			method.ImplName = x.ShortName
			noImplFunc = append(noImplFunc, method)
			continue
		}
		if !method.Equal(sMth) {
			equal = false
			changeFunc = append(changeFunc, method)
			continue
		}
	}
	return equal, noImplFunc, changeFunc
}

type XMethod struct {
	ImplName   string
	Name       string
	Params     []XArg
	Results    []XArg
	Comment    string
	Sort       int
	HTTPRule   string
	HTTPMethod string
}

func (x XMethod) Equal(x2 XMethod) bool {
	if x.Name != x2.Name {
		return false
	}
	if len(x.Params) != len(x2.Params) {
		return false
	}
	for i, it := range x.Params {
		if it.Type != x2.Params[i].Type {
			return false
		}
	}
	if len(x.Results) != len(x2.Results) {
		return false
	}
	for i, it := range x.Results {
		if it.Type != x2.Results[i].Type {
			return false
		}
	}
	return true
}

type XField struct {
	Name    string
	Type    string
	SType   int
	Idx     int
	Tag     string
	Comment string
}

func (x XField) Equal(x2 XField) bool {
	if x.Type != x2.Type {
		return false
	}
	if x.Name != x2.Name {
		return false
	}
	if x.Tag != x2.Tag {
		return false
	}
	if x.Comment != x2.Comment {
		return false
	}
	return true
}

type XArg struct {
	Name string
	Type string
}

type TagDesc struct {
	JSON     *TagItem
	PB       *TagItem
	DB       *TagItem
	Validate *TagItem
}

type TagItem struct {
	Name string            `json:"name"`
	Txt  string            `json:"txt"`
	Opts map[string]string `json:"opts"`
}

func (x XField) GetTag(tag string) *TagItem {
	re, _ := regexp.Compile(tag + `:"(\S+)"`)
	if rs := re.FindStringSubmatch(x.Tag); len(rs) > 0 {
		txt := rs[1]
		if txt == "-" {
			return &TagItem{
				Name: "-",
				Txt:  "-",
				Opts: map[string]string{},
			}
		}
		it := &TagItem{
			Name: txt,
			Txt:  txt,
			Opts: map[string]string{},
		}
		if strings.Contains(txt, ":") || strings.Contains(txt, ";") {
			it.Name = x.Name
			tmp := strings.Split(txt, ";")
			for idx, s := range tmp {
				if idx == 0 {
					it.Name = s
				}
				r := strings.Split(s, ":")
				if len(r) == 2 {
					it.Opts[r[0]] = r[1]
					if r[0] == "column" {
						it.Name = r[1]
					}
				} else {
					it.Opts[r[0]] = ""
				}
			}
		}

		return it
	}
	return nil
}

type EntryModule struct {
	Name       string
	Key        string
	Middleware []string
	WithCommon bool
	FuncList   []*EntryModuleFunc
}

type EntryModuleFunc struct {
	idx        int
	Name       string
	Key        string
	KeyLi      string
	Middleware []string
	Request    *EntryModuleFuncReq
	Response   *EntryModuleFuncResp
}

type EntryModuleFuncReq struct {
	Name string `json:"name"`
	XST  XST    `json:"xst"`
}

type EntryModuleFuncResp struct {
	Name string `json:"name"`
	XST  XST    `json:"xst"`
}
