package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DB struct {
	conn *sqlx.DB
}

type Column struct {
	ColumnName    string `db:"column_name"`
	DataType      string `db:"data_type"`
	IsNullable    string `db:"is_nullable"`
	TableName     string `db:"table_name"`
	ColumnComment string `db:"column_comment"`
	ColumnKey     string `db:"column_key"`
}

func GetDB(dsn string) *DB {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return &DB{conn: db}
}

func (db *DB) TableCreateSQL(tableName string) (str string) {
	var t1 string
	row := db.conn.QueryRowx(fmt.Sprintf("show create table %s;", tableName))
	err := row.Scan(&t1, &str)
	if err != nil {
		return ""
	}
	return str
}

func (db *DB) TableColumns(tableName string) (columnMap map[string]Column, list []Column) {
	tmp := strings.Split(tableName, ".")
	tableName = tmp[len(tmp)-1]
	var columns = make([]Column, 0)
	err := db.conn.Select(&columns, `SELECT column_name,data_type,is_nullable,table_name,column_comment,column_key
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE() AND table_name = ?`, tableName)
	if err != nil {
		return
	}
	columnMap = make(map[string]Column, 0)
	for _, it := range columns {
		columnMap[it.ColumnName] = it
	}
	return columnMap, columns
}

func AddStrSqlC(str string, cs string) string {
	tmp := strings.Split(str, ".")
	tmp2 := []string{}
	for _, s := range tmp {
		tmp2 = append(tmp2, fmt.Sprintf("%s%s%s", cs, s, cs))
	}
	return strings.Join(tmp2, ".")
}
