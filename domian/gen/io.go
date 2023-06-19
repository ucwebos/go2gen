package gen

import (
	"fmt"
	"go2gen/conf"
	"log"
	"sort"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"
	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
)

func (m *Manager) IOApi(xsts map[string]parser.XST) {
	var (
		buf  = []byte{}
		buf2 = []byte{}
	)
	for _, xst := range xsts {
		b1, b2, err := m.IO(xst)
		if err != nil {
			log.Panicf("gen io err: %v \n", err)
		}
		buf = append(buf, b1...)
		buf2 = append(buf2, b2...)
	}

	filename := fmt.Sprintf("%s/type_gen.api", m.AppPath)
	buf = append([]byte("// Code generated by go2gen. DO NOT EDIT.\n"), buf...)
	buf = m.format(buf, filename)
	err := tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("io gen [%s] write file err: %v \n", filename, err)
	}

	filename2 := fmt.Sprintf("%s/io_converter_gen.go", m.Tmpl.ConvIODir)

	bufH := m.GenFileHeader("converter", []string{
		fmt.Sprintf("%s/common/tools/tool_time", conf.Global.ProjectName),
		fmt.Sprintf("%s/%s/internal/domain/entity", conf.Global.ProjectName, m.AppName),
		fmt.Sprintf("%s/%s/internal/types", conf.Global.ProjectName, m.AppName),
	})
	buf2 = append(bufH, buf2...)
	buf2 = m.format(buf2, filename)
	err = tool_file.WriteFile(filename2, buf2)
	if err != nil {
		log.Printf("io conv gen [%s] write file err: %v \n", filename, err)
	}
}

func (m *Manager) IO(xst parser.XST) ([]byte, []byte, error) {
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
		tagIO := field.GetTag("io")
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

		//type2 := ""
		tags := fmt.Sprintf("`json:\"%s\"`", tagJSON.Name)
		fType := field.Type
		switch field.SType {
		case 1:
			//type2 = strings.Replace(field.Type, "*", "", 1)
			if strings.Contains(field.Type, "time.Time") {
				//type2 = strings.Replace(field.Type, "*", "", 1)
				field.SType = parser.STypeTime
				fType = "string"
			}
		case 2:
			//type2 = strings.Replace(field.Type, "[]", "", 1)
			//fType = AddEntityPkg(fType)
		case 3:
			//fType = AddEntityPkg(fType)
		}

		gio.Fields = append(gio.Fields, tpls.IoField{
			Name: field.Name,
			Type: fType,
			//Type2:   type2,
			SType:   field.SType,
			Tag:     tags,
			Comment: field.Comment,
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
