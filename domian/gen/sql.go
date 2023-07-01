package gen

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"go2gen/domian/gen/tpls"
	"go2gen/domian/parser"
	"go2gen/utils"

	"github.com/xbitgo/core/tools/tool_file"
)

func (m *Manager) Do2Sql(dsn string) error {
	db := utils.GetDB(dsn)
	ipr, err := parser.Scan(m.Tmpl.DoDir, parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("do2Sql: parse dir[%s], err: %v", m.Tmpl.DoDir, err)
	}
	for s, xst := range ipr.StructList {
		if v, ok := ipr.ConstStrList["TableName"+s]; ok {
			err = m.createTableSQL(db, v, xst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) createTableSQL(db *utils.DB, tableName string, xst parser.XST) error {
	filename := fmt.Sprintf("%s/%s_create.sql", m.Tmpl.SQLDir, tableName)
	createSQL := db.TableCreateSQL(tableName)
	if createSQL != "" {
		tool_file.WriteFile(filename, []byte(createSQL))
		m.modifySQL(db, tableName, xst)
	} else {
		genSql := m.toGenSQL(tableName, xst)
		createSQL, err := genSql.CreateTable()
		if err != nil {
			return err
		}
		tool_file.WriteFile(filename, createSQL)
	}
	return nil
}

func (m *Manager) toGenSQL(tableName string, xst parser.XST) tpls.GenSQL {
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	primaryKey := ""
	last := ""
	sFields := make([]tpls.SQLField, 0)
	for _, field := range fieldList {
		dbTag := field.GetTag("db")
		if dbTag != nil {
			name := utils.AddStrSqlC(dbTag.Name, "`")
			if sf, ok := tpls.SpecialField[dbTag.Name]; ok {
				sf.TableName = utils.AddStrSqlC(tableName, "`")
				sf.After = last
				sf.SrcName = dbTag.Name
				sFields = append(sFields, sf)
				last = name
				continue
			}
			if strings.Contains(field.Tag, "primaryKey") {
				primaryKey = name
			}
			tt := tpls.TypeMap[strings.TrimPrefix(field.Type, "*")]
			sf := tpls.SQLField{
				TableName: utils.AddStrSqlC(tableName, "`"),
				After:     last,
				Name:      name,
				SrcName:   dbTag.Name,
				Type:      tt.Type,
				DataType:  tt.DataType,
				Default:   tt.Default,
				Comment:   field.Comment,
				NotNull:   "NOT NULL",
			}
			sFields = append(sFields, sf)
			last = name
		}
	}
	genSql := tpls.GenSQL{
		TableName:  utils.AddStrSqlC(tableName, "`"),
		PrimaryKey: primaryKey,
		Fields:     sFields,
	}
	return genSql
}

func (m *Manager) modifySQL(db *utils.DB, tableName string, xst parser.XST) error {
	columns, _ := db.TableColumns(tableName)
	genSql := m.toGenSQL(tableName, xst)
	addColumns := make([]tpls.SQLField, 0)
	for _, field := range genSql.Fields {
		if _, ok := columns[field.SrcName]; !ok {
			addColumns = append(addColumns, field)
		}
	}
	if len(addColumns) > 0 {
		filename := fmt.Sprintf("%s/%s_column_add_%s.sql", m.Tmpl.SQLDir, tableName, time.Now().Format("200601021504"))
		genSql.Fields = addColumns
		createSQL, err := genSql.AddColumns()
		if err != nil {
			fmt.Println(err)
			return err
		}
		tool_file.WriteFile(filename, createSQL)
	}
	return nil
}
