package gen

import (
	"fmt"
	"go2gen/conf"
	"log"
	"sort"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
)

func (m *Manager) Do(xst parser.XST) error {
	gdo := tpls.Do{
		SrcPath: "",
		Package: "do",
		Imports: make([]string, 0),
		Name:    xst.Name,
		Fields:  make([]tpls.DoField, 0),
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	for _, field := range fieldList {
		tagDesc := field.GetTag("db")
		if tagDesc != nil {
			tag := tagDesc.Txt
			convSlice := false
			isPoint := false
			type2 := ""
			if tagDesc.Opts != nil && len(tagDesc.Opts) > 0 {
				if v, ok := tagDesc.Opts["conv"]; ok {
					tagConv := fmt.Sprintf("conv:%s", v)
					convSlice = true
					tag = strings.Replace(tag, tagConv+";", "", 1)
					tag = strings.Replace(tag, tagConv, "", 1)
				}
			}
			tags := fmt.Sprintf("`db:\"%s\" gorm:\"%s\"`", tagDesc.Name, tag)
			fType := field.Type
			switch field.SType {
			case 1:
				type2 = strings.Replace(field.Type, "*", "", 1)
				if strings.Contains(field.Type, "time.Time") {
					field.SType = parser.STypeTime
				} else {
					if strings.Index(field.Type, "*") == 0 {
						isPoint = true
					}
				}

			case 2:
				type2 = strings.Replace(field.Type, "[]", "", 1)
				if strings.Contains(type2, "[]") || strings.Index(type2, ".") > 0 {
					convSlice = false
				}
				if tool_str.UFirst(type2) || strings.Contains(type2, "map") {
					convSlice = false
				}
				fType = AddEntityPkg(fType)
			case 3:
				fType = AddEntityPkg(fType)
			}

			gdo.Fields = append(gdo.Fields, tpls.DoField{
				Name:      field.Name,
				Type:      fType,
				Type2:     type2,
				SType:     field.SType,
				Tag:       tags,
				ConvSlice: convSlice,
				IsPoint:   isPoint,
				Comment:   field.Comment,
			})
		}
	}
	if len(gdo.Fields) == 0 {
		return nil
	}
	filename := fmt.Sprintf("%s/%s_do_gen.go", m.Tmpl.DoDir, tool_str.ToSnakeCase(xst.Name))
	buf, err := gdo.Execute()
	if err != nil {
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("do gen [%s] write file err: %v \n", filename, err)
	}
	m.dbConv(xst, gdo)

	return nil
}

func (m *Manager) DoTypeDef() {
	ipr, err := parser.Scan(m.Tmpl.DoDir, parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("do2Sql: parse dir[%s], err: %v", m.Tmpl.DoDir, err)
	}
	bufs := []byte(fmt.Sprintf(tpls.EntityTypeDefCodes, "do"))
	for _, xst := range ipr.StructList {
		buf, err := m._typedef(xst)
		if err != nil {
			log.Printf("gen mapType err: %v \n", err)
		}
		bufs = append(bufs, buf...)
	}
	filename := fmt.Sprintf("%s/typedef_code_gen.go", m.Tmpl.DoDir)
	bufs = m.format(bufs, filename)
	err = tool_file.WriteFile(filename, bufs)
	if err != nil {
		log.Printf("typedef gen [%s] write file err: %v \n", filename, err)
	}
}

func (m *Manager) dbConv(xst parser.XST, gdo tpls.Do) error {
	convGen := tpls.DoConv{
		Name:    gdo.Name,
		Package: "converter",
		Imports: []string{
			"encoding/json",
			fmt.Sprintf("%s/%s/internal/domain/entity", conf.Global.ProjectName, m.AppName),
			fmt.Sprintf("%s/%s/internal/domain/repo/dbal/do", conf.Global.ProjectName, m.AppName),
		},
		Fields: gdo.Fields,
	}
	filename := fmt.Sprintf("%s/%s_converter_gen.go", m.Tmpl.ConvDoDir, tool_str.ToSnakeCase(xst.Name))
	buf, err := convGen.Execute()
	if err != nil {
		return err
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("do gen [%s] write file err: %v \n", filename, err)
	}

	return nil
}
