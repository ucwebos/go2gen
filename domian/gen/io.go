package gen

import (
	"fmt"
	"github.com/ucwebos/go2gen/conf"
	"github.com/xbitgo/core/tools/tool_str"
	"log"
	"sort"
	"strings"

	"github.com/ucwebos/go2gen/domian/gen/tpls"
	"github.com/ucwebos/go2gen/domian/parser"
	"github.com/xbitgo/core/tools/tool_file"
)

func (m *Manager) IOEntries(xsts map[string]parser.XST, entry string) {
	var (
		buf  = []byte{}
		buf2 = []byte{}
	)

	for _, xst := range xsts {
		b1, b2, err := m.IO(xst, entry)
		if err != nil {
			log.Panicf("gen io err: %v \n", err)
		}
		buf = append(buf, b1...)
		buf2 = append(buf2, b2...)
	}

	filename := fmt.Sprintf("%s/internal/entry/%s/types/types_gen.go", m.AppPath, entry)
	buf = append([]byte("// Code generated by go2gen. DO NOT EDIT.\n"+
		"package types \n"), buf...)
	buf = m.format(buf, filename)
	err := tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("io gen [%s] write file err: %v \n", filename, err)
	}

	filename2 := fmt.Sprintf("%s/%s/converter/io_converter_gen.go", m.Tmpl.EntryDir, entry)

	bufH := m.GenFileHeader("converter", []string{
		fmt.Sprintf("%s/common/tools/tool_time", conf.Global.ProjectName),
		fmt.Sprintf("%s/%s/internal/domain/entity", conf.Global.ProjectName, m.AppName),
		fmt.Sprintf("%s/%s/internal/entry/%s/types", conf.Global.ProjectName, m.AppName, entry),
	})
	buf2 = append(bufH, buf2...)
	buf2 = m.format(buf2, filename)
	err = tool_file.WriteFile(filename2, buf2)
	if err != nil {
		log.Printf("io conv gen [%s] write file err: %v \n", filename, err)
	}
}

func (m *Manager) IO(xst parser.XST, tagName string) ([]byte, []byte, error) {
	gio := tpls.IO{
		Name:   xst.Name,
		Fields: make([]tpls.IoField, 0),
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	for _, field := range fieldList {
		tagJSON := field.GetTag("json")
		tagIO := field.GetTag(tagName)
		if tagJSON == nil {
			continue
		}
		if tagIO != nil {
			if tagIO.Txt == "-" {
				continue
			}
			if tagIO.Txt != "" {
				tagJSON.Name = tagIO.Name
			}
		}

		type2 := ""
		type2Entity := false
		tags := fmt.Sprintf("`json:\"%s\"`", tagJSON.Name)
		fType := field.Type
		switch field.SType {
		case 1:
			type2 = strings.Replace(field.Type, "*", "", 1)
			if strings.Contains(field.Type, "time.Time") {
				//type2 = strings.Replace(field.Type, "*", "", 1)
				field.SType = parser.STypeTime
				fType = "string"
			}
		case 2:
			type2 = strings.Replace(field.Type, "[]", "", 1)
			type2 = strings.Replace(type2, "*", "", 1)
			if tool_str.UFirst(type2) {
				type2Entity = true
			}
			//fType = AddEntityPkg(fType)
		case 3:
			//fType = AddEntityPkg(fType)
		}

		gio.Fields = append(gio.Fields, tpls.IoField{
			Name:        field.Name,
			Type:        fType,
			Type2:       type2,
			Type2Entity: type2Entity,
			SType:       field.SType,
			Tag:         tags,
			Comment:     field.Comment,
		})
	}

	if len(gio.Fields) == 0 {
		return nil, nil, nil
	}
	buf, err := gio.Execute()
	if err != nil {
		return nil, nil, err
	}
	convBuf, err := m.ioConv(xst, gio)
	if err != nil {
		return nil, nil, err
	}
	return buf, convBuf, nil
}

func (m *Manager) ioConv(xst parser.XST, gio tpls.IO) ([]byte, error) {
	convGen := tpls.IoConv{
		Name:   gio.Name,
		Fields: gio.Fields,
	}
	buf, err := convGen.Execute()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (m *Manager) GenFileHeader(pkg string, imports []string) []byte {
	importStr := ""
	for _, i := range imports {
		importStr += fmt.Sprintf(`"%s"`+"\n", i)
	}
	bufH := fmt.Sprintf(`// Code generated by go2gen. DO NOT EDIT.
package %s

import (
%s
)`, pkg, importStr)
	return []byte(bufH)
}
