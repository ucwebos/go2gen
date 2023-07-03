package gen

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/xbitgo/core/tools/tool_file"
	"go2gen/conf"
	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
	"go2gen/utils"
	"log"
	"os"
	"strings"
)

func (m *Manager) HandlerAndDoc(entry string) error {
	ips := &parser.IParser{
		Pwd:          fmt.Sprintf("%s/%s/types", m.Tmpl.EntryDir, entry),
		Package:      "types",
		INFList:      make(map[string]parser.INF, 0),
		StructList:   make(map[string]parser.XST, 0),
		OtherStruct:  make(map[string]parser.XST, 0),
		ConstStrList: make(map[string]string),
		BindFuncMap:  map[string]map[string]parser.XMethod{},
		NewFuncList:  map[string]parser.XMethod{},
		EntryModules: []parser.EntryModule{},
		ParseType:    parser.ParseTypeHandler,
	}
	err := ips.ParseFile(fmt.Sprintf("%s/types_io.go", ips.Pwd))
	if err != nil {
		return err
	}
	// handler
	dir := fmt.Sprintf("%s/%s/handler", m.Tmpl.EntryDir, entry)
	ips2, err := parser.Scan(dir, parser.ParseTypeImpl)
	if err != nil {
		return err
	}
	for _, module := range ips.EntryModules {
		m.handler(entry, ips2.FuncList, module)
	}
	m.routes(entry, ips.EntryModules)

	m.docs(entry, ips.EntryModules)

	return nil
}

func (m *Manager) routes(entry string, modules []parser.EntryModule) error {
	filename := fmt.Sprintf("%s/%s/routes_gen.go", m.Tmpl.EntryDir, entry)
	t := &tpls.Routes{
		ProjectName: conf.Global.ProjectName,
		AppName:     m.AppName,
		Entry:       entry,
		ModuleList:  modules,
	}
	buf, err := t.Execute()
	if err != nil {
		log.Printf("routes gen Execute err: %v \n", err)
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("routes gen [%s] write file err: %v \n", filename, err)
	}
	return err

}

func (m *Manager) handler(entry string, hasFuncMap map[string]parser.XMethod, module parser.EntryModule) error {
	filename := fmt.Sprintf("%s/%s/handler/%s.go", m.Tmpl.EntryDir, entry, module.Key)
	//fmt.Println(filename)
	var (
		buf []byte
		err error
	)
	if tool_file.Exists(filename) {
		buf, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
		t := &tpls.HandlerFuncAppend{
			Body:     buf,
			FuncList: make([]*parser.EntryModuleFunc, 0),
		}
		for _, it := range module.FuncList {
			if _, ok := hasFuncMap[it.Key]; !ok {
				t.FuncList = append(t.FuncList, it)
			}
		}
		buf, err = t.Execute()
	} else {
		t := &tpls.HandlerFunc{
			ProjectName: conf.Global.ProjectName,
			AppName:     m.AppName,
			Entry:       entry,
			FuncList:    module.FuncList,
		}
		buf, err = t.Execute()

	}
	if err != nil {
		log.Printf("handler gen [%s] write file err: %v \n", filename, err)
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("handler gen [%s] write file err: %v \n", filename, err)
	}
	return err
}

func (m *Manager) docs(entry string, modules []parser.EntryModule) error {
	ips := &parser.IParser{
		Pwd:          fmt.Sprintf("%s/%s/types", m.Tmpl.EntryDir, entry),
		Package:      "types",
		INFList:      make(map[string]parser.INF, 0),
		StructList:   make(map[string]parser.XST, 0),
		OtherStruct:  make(map[string]parser.XST, 0),
		ConstStrList: make(map[string]string),
		BindFuncMap:  map[string]map[string]parser.XMethod{},
		NewFuncList:  map[string]parser.XMethod{},
		EntryModules: []parser.EntryModule{},
		ParseType:    parser.ParseTypeWatch,
	}
	err := ips.ParseFile(fmt.Sprintf("%s/types_gen.go", ips.Pwd))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//
	for _, module := range modules {
		dir := fmt.Sprintf("%s/docs/%s/%s", m.AppPath, entry, module.Key)
		os.MkdirAll(dir, 0777)
		for _, moduleFunc := range module.FuncList {
			m.docsItem(entry, module.Key, dir, moduleFunc, ips.StructList)
		}
		//
	}
	filename := fmt.Sprintf("%s/docs/%s/_sidebar.md", m.AppPath, entry)
	sider, err := os.ReadFile(filename)
	siderStr := string(sider)
	idx1 := strings.Index(siderStr, "---")
	idx2 := strings.LastIndex(siderStr, "---")

	t := &tpls.DocsSidebar{
		Entry:   entry,
		Modules: modules,
	}
	buf, err := t.Execute()
	_str := siderStr[:idx1+3] + string(buf) + siderStr[idx2:]
	tool_file.WriteFile(filename, []byte(_str))
	return nil
}

func (m *Manager) docsItem(entry, moduleKey, dir string, f *parser.EntryModuleFunc, structList map[string]parser.XST) {

	filename := fmt.Sprintf("%s/%s.md", dir, f.KeyLi)
	request := make([]tpls.DocsItemField, 0)
	if f.Request != nil {
		request = append(request, m.toDocsItemFields(f.Request.XST.FieldList, structList, "")...)
	}
	response := make([]tpls.DocsItemField, 0)
	if f.Response != nil {
		response = append(response, m.toDocsItemFields(f.Response.XST.FieldList, structList, "")...)
	}
	t := &tpls.DocsItem{
		Name:      f.Name,
		RoutePath: fmt.Sprintf("/%s/%s/%s", entry, moduleKey, f.KeyLi),
		Request:   request,
		Response:  response,
		ExpJSON:   []byte{},
	}
	if f.Response != nil {
		body := m.getJSON(f.Response.XST.FieldList, structList)
		sb, _ := jsoniter.MarshalIndent(body, "", "  ")
		t.ExpJSON = append(t.ExpJSON, []byte("```\n")...)
		t.ExpJSON = append(t.ExpJSON, sb...)
		t.ExpJSON = append(t.ExpJSON, []byte("\n```")...)

	}
	buf, err := t.Execute()
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("doc item gen [%s] write file err: %v \n", filename, err)
	}
	return
}

func (m *Manager) getJSON(fields map[string]parser.XField, structList map[string]parser.XST) *utils.OrderMap {
	body := make(map[string]interface{})
	for _, it := range fields {
		j := it.GetTag("json")
		body[j.Name] = m.getJSONVal(it, structList)
	}
	om := utils.NewOrderMap(utils.DefaultOrderMapKeySort)
	_ = om.LoadStringMap(body)
	return om
}

func (m *Manager) getJSONVal(field parser.XField, structList map[string]parser.XST) interface{} {
	switch field.SType {
	case parser.STypeStruct:
		sk := strings.TrimPrefix(field.Type, "*")
		if v, ok := structList[sk]; ok {
			return m.getJSON(v.FieldList, structList)
		}
	case parser.STypeSlice:
		sk := strings.ReplaceAll(strings.ReplaceAll(field.Type, "*", ""), "[]", "")
		if v, ok := structList[sk]; ok {
			return []interface{}{
				m.getJSON(v.FieldList, structList),
			}
		} else {
			return []interface{}{m.getZeroVal(field.Type)}
		}
	default:
		return m.getZeroVal(field.Type)
	}
	return ""
}
func (m *Manager) getZeroVal(xType string) interface{} {
	if xType == "bool" {
		return true
	}
	if strings.Contains(xType, "int") {
		return 0
	}
	if strings.Contains(xType, "float") {
		return 0.1
	}
	return ""
}

func (m *Manager) toDocsItemFields(fields map[string]parser.XField, structList map[string]parser.XST, prefix string) []tpls.DocsItemField {
	_fields := m.sortFields(fields)
	request := make([]tpls.DocsItemField, 0)
	for _, field := range _fields {
		j := field.GetTag("json")
		name := prefix + j.Name
		_type := field.Type
		switch field.SType {
		case parser.STypeStruct:
			prefix = strings.ReplaceAll(prefix, "[i].", "")
			_type = "object"
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
			sk := strings.TrimPrefix(field.Type, "*")
			if v, ok := structList[sk]; ok {
				r := m.toDocsItemFields(v.FieldList, structList, prefix+"&emsp;&emsp;")
				request = append(request, r...)
			}
		case parser.STypeSlice:
			prefix = strings.ReplaceAll(prefix, "[i].", "")
			_type = "array"
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
			sk := strings.ReplaceAll(strings.ReplaceAll(field.Type, "*", ""), "[]", "")
			if v, ok := structList[sk]; ok {
				r := m.toDocsItemFields(v.FieldList, structList, prefix+"&emsp;&emsp;[i].")
				request = append(request, r...)
			}
		default:
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
		}

	}
	return request
}

func (m *Manager) sortFields(fields map[string]parser.XField) []parser.XField {
	r := make([]parser.XField, len(fields))
	for _, field := range fields {
		r[field.Idx] = field
	}
	return r

}
