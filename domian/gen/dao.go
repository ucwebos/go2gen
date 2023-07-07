package gen

import (
	"fmt"
	"github.com/ucwebos/go2gen/conf"
	"log"
	"strings"

	"github.com/ucwebos/go2gen/domian/gen/tpls"
	"github.com/ucwebos/go2gen/domian/parser"
	"github.com/xbitgo/core/tools/tool_file"
)

func (m *Manager) Dao() error {
	ipr, err := parser.Scan(m.Tmpl.DoDir, parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("do2Sql: parse dir[%s], err: %v", m.Tmpl.DoDir, err)
	}
	buf := []byte{}
	for _, xst := range ipr.StructList {
		b, err := m.genDao(xst)
		if err != nil {
			log.Panicf("gen dao err: %v", err)
		}
		buf = append(buf, b...)
	}

	filename := fmt.Sprintf("%s/dao_gen.go", m.Tmpl.DaoDir)
	bufH := m.GenFileHeader("dao", []string{
		"github.com/pkg/errors",
		"gorm.io/gorm",
		fmt.Sprintf("%s/%s/internal/domain/repo/dbal/do", conf.Global.ProjectName, m.AppName),
	})
	buf = append(bufH, buf...)
	buf = m.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("dao gen [%s] write file err: %v \n", filename, err)
	}

	return nil
}

func (m *Manager) genDao(xst parser.XST) ([]byte, error) {
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
		EntityName:     xst.Name,
		DaoName:        strings.TrimSuffix(xst.Name, "Do") + "Dao",
		EntityListName: fmt.Sprintf("%sList", xst.Name),
		TableName:      fmt.Sprintf("do.TableName%s", xst.Name),
		PkName:         pkName,
		PkType:         pkType,
		PkCol:          pkCol,
	}

	//filename := fmt.Sprintf("%s/%s_dao_gen.go", m.Tmpl.DaoDir, tool_str.ToSnakeCase(strings.TrimSuffix(xst.Name, "Do")))
	buf, err := dao.Execute()
	if err != nil {
		return buf, err
	}
	return buf, nil
	//buf = m.format(buf, filename)
	//err = tool_file.WriteFile(filename, buf)
	//if err != nil {
	//	log.Printf("do gen [%s] write file err: %v \n", filename, err)
	//}
}
