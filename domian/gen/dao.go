package gen

import (
	"fmt"
	"go2gen/conf"
	"log"
	"strings"

	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"

	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
)

func (m *Manager) Dao() error {
	ipr, err := parser.Scan(m.Tmpl.DoDir, parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("do2Sql: parse dir[%s], err: %v", m.Tmpl.DoDir, err)
	}
	for _, xst := range ipr.StructList {
		m.genDao(xst)
	}
	return nil
}

func (m *Manager) genDao(xst parser.XST) {
	var (
		pkName = ""
		pkType = ""
		pkCol  = ""
	)
	for _, field := range xst.FieldList {
		tag := field.GetTag("gorm")
		if tag != nil && strings.Contains(tag.Txt, "primaryKey") {
			pkName = field.Name
			pkType = field.Type
			pkCol = tag.Name
		}
	}
	dao := tpls.Dao{
		EntityPackage:  fmt.Sprintf("%s/%s/internal/domain/repo/dbal/do", conf.Global.ProjectName, m.AppName),
		EntityName:     xst.Name,
		DaoName:        strings.TrimSuffix(xst.Name, "Do") + "Dao",
		EntityListName: fmt.Sprintf("%sList", xst.Name),
		TableName:      fmt.Sprintf("do.TableName%s", xst.Name),
		PkName:         pkName,
		PkType:         pkType,
		PkCol:          pkCol,
	}

	filename := fmt.Sprintf("%s/%s_dao_gen.go", m.Tmpl.DaoDir, tool_str.ToSnakeCase(strings.TrimSuffix(xst.Name, "Do")))
	buf, err := dao.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("do gen [%s] write file err: %v \n", filename, err)
	}
}
