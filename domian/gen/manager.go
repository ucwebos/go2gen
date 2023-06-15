package gen

import (
	"fmt"
	"go/format"
	"go2gen/conf"
	"log"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/xbitgo/core/di"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
)

type projectINF interface {
	RootPath() string
}

type Manager struct {
	Tmpl    *conf.Tmpl
	AppName string
	AppPath string
	Project projectINF `di:"project"`
}

func NewManager(tmpl *conf.Tmpl, appName string, appPath string) *Manager {
	m := &Manager{
		Tmpl:    tmpl,
		AppName: appName,
		AppPath: appPath,
	}
	di.MustBind(m)
	return m
}

func (m *Manager) EntityTypeDef() {
	ipr, err := parser.Scan(m.Tmpl.EntityDir, parser.ParseTypeWatch)
	if err != nil {
		log.Fatalf("EntityTypeDef: parse dir[%s], err: %v", m.Tmpl.EntityDir, err)
	}
	bufs := []byte(fmt.Sprintf(tpls.EntityTypeDefCodes, "entity"))
	for _, xst := range ipr.StructList {
		buf, err := m._typedef(xst)
		if err != nil {
			log.Printf("gen mapType err: %v \n", err)
		}
		bufs = append(bufs, buf...)
	}
	filename := fmt.Sprintf("%s/typedef_code_gen.go", m.Tmpl.EntityDir)
	bufs = m.format(bufs, filename)
	err = tool_file.WriteFile(filename, bufs)
	if err != nil {
		log.Printf("do gen [%s] write file err: %v \n", filename, err)
	}
}

func (m *Manager) _typedef(xst parser.XST) ([]byte, error) {
	tGen := tpls.EntityTypeMap{
		EntityName:     xst.Name,
		EntityListName: fmt.Sprintf("%sList", xst.Name),
		Field:          make([]tpls.Field, 0),
		HasCreator:     false,
		CreatorName:    "",
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	feList := make([]tpls.Field, 0)

	for _, field := range fieldList {
		_type := field.Type
		tags := strings.Trim(field.Tag, "`")
		tagsMap := parseFieldTagMap(tags)
		dbTag := tagsMap["db"]
		if dbTag != "" && strings.Contains(dbTag, ";") {
			dbTag = strings.Split(dbTag, ";")[0]
		}
		if dbTag == "create_time" || dbTag == "update_time" || dbTag == "id" || dbTag == "deleted_at" {
			dbTag = ""
		}
		fe := tpls.Field{
			Field:           field.Name,
			FieldTag:        tags,
			FieldEscapedTag: fmt.Sprintf("%q", tags),
			FieldTagMap:     tagsMap,
			DBTag:           dbTag,
			Type:            _type,
			UseJSON:         false,
			NamedType:       "",
			TypeInName:      "",
			GenSliceFunc:    true,
			Nullable:        false,
			Comparable:      false,
		}
		if field.SType != 0 && field.SType != 4 {
			fe.UseJSON = true
		}
		if strings.Index(_type, "*") == 0 || field.SType >= 2 || _type == "interface{}" {
			fe.Nullable = true
		} else {
			fe.Comparable = true
		}
		switch _type {
		case "int":
			fe.TypeInName = "Int"
		case "int32":
			fe.TypeInName = "Int32"
		case "int64":
			fe.TypeInName = "Int64"
		case "string":
			fe.TypeInName = "String"
		default:
			fe.GenSliceFunc = false
		}
		feList = append(feList, fe)
	}

	tGen.Field = feList
	return tGen.Execute()
}

func AddEntityPkg(str string) string {
	for i, r := range []rune(str) {
		if tool_str.UFirst(string(r)) {
			return str[:i] + "entity." + str[i:]
		}
	}
	return str
}

func (m *Manager) ConfGI(pkg string, xst parser.XST) error {
	//gSdi := tpls.NewSDI(pkg, xst)
	//filename := fmt.Sprintf("%s/%s", m.Tmpl.ConfDir, "di_register_gen.go")
	//buf, err := gSdi.Execute()
	//if err != nil {
	//	return err
	//}
	//buf = m.format(buf, filename)
	//err = tool_file.WriteFile(filename, buf)
	//if err != nil {
	//	log.Printf("app gen [%s] write file err: %v \n", filename, err)
	//}

	return nil
}

func (m *Manager) GI(iParser *parser.IParser) error {
	gdi := tpls.DI{
		Pkg:  iParser.Package,
		List: make([]tpls.DItem, 0),
	}
	for _, xst := range iParser.StructList {
		if xst.GI {
			it := tpls.DItem{
				Name:    xst.Name,
				NameVal: tool_str.ToLFirst(xst.Name),
			}
			if nMth, ok := iParser.NewFuncList[xst.Name]; ok {
				it.NewReturnsLen = len(nMth.Results)
			}
			gdi.List = append(gdi.List, it)
		}
	}
	if len(gdi.List) == 0 {
		return nil
	}
	filename := fmt.Sprintf("%s/%s", iParser.Pwd, "gi_gen.go")
	buf, err := gdi.Execute()
	if err != nil {
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("app gen [%s] write file err: %v \n", filename, err)
	}
	return nil
}

func (m *Manager) Protoc(pbFile string) error {
	cmd := exec.Command("protoc", "-I", ".", "-I", "./third_party", "--gogofast_out", "../../", "--go-grpc_out", "../../", "--swagger_out=logtostderr=true:.", pbFile)
	if strings.HasSuffix(pbFile, "_gen.proto") || strings.HasSuffix(pbFile, "/base.proto") {
		cmd = exec.Command("protoc", "-I", ".", "-I", "./third_party", "--gogofast_out", "../../", "--go-grpc_out", "../../", pbFile)
	}
	cmd.Dir = m.Project.RootPath() + "/proto"
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(pbFile)
		log.Printf("run main error %v", err)
	}
	return nil
}

func (m *Manager) format(buf []byte, filename string) []byte {
	buf2, err := format.Source(buf)
	if err == nil {

		buf = buf2
	}
	buf3, err := imports.Process(filename, buf, nil)
	if err == nil {
		buf = buf3
	}
	return buf
}

func parseFieldTagMap(tag string) map[string]string {
	tagParts := strings.Split(tag, "\" ")
	result := map[string]string{}
	for _, parts := range tagParts {
		pairs := strings.SplitN(parts, ":", 2)
		if len(parts) < 2 {
			continue
		}
		result[pairs[0]] = strings.Trim(pairs[1], "\"")
	}
	return result
}
