package tpls

import (
	"bytes"
	"text/template"
)

const sqlCreateTableTpl = `
-- 注意！此为根据代码生成的SQL需要根据实际情况人工修改字段类型/长度等
CREATE TABLE {{.TableName}} (
{{- range .Fields}}
	{{.Name}} {{.Type}} {{.NotNull}} {{.Default}} COMMENT '{{.Comment}}',
{{- end}}
	PRIMARY KEY ({{.PrimaryKey}}),
	KEY (` + "`" + `deleted_at` + "`" + `)
) ENGINE={{.Engine}} DEFAULT CHARSET={{.Charset}} COLLATE={{.Collate}};
`

const sqlCreateAddColumns = `
-- 注意！此为根据代码生成的SQL需要根据实际情况人工修改字段类型/长度等
{{- range .Fields}}
ALTER TABLE {{.TableName}}  
ADD COLUMN {{.Name}} {{.Type}} {{.NotNull}} {{.Default}} COMMENT '{{.Comment}}' AFTER {{.After}};
{{- end}}
`

var (
	TypeMap = map[string]SQLField{
		"int": {
			Type:     "bigint(20)",
			DataType: "bigint",
			Default:  "DEFAULT 0",
		},
		"int64": {
			Type:     "bigint(20)",
			DataType: "bigint",
			Default:  "DEFAULT 0",
		},
		"int32": {
			Type:     "int(11)",
			DataType: "int",
			Default:  "DEFAULT 0",
		},
		"int16": {
			Type:     "int(9)",
			DataType: "int",
			Default:  "DEFAULT 0",
		},
		"int8": {
			Type:     "tinyint(4)",
			DataType: "tinyint",
			Default:  "DEFAULT 0",
		},
		"string": {
			Type:     "varchar(255)",
			DataType: "varchar",
			Default:  "DEFAULT ''",
		},
		"float64": {
			Type:     "decimal(10, 3)",
			DataType: "decimal",
			Default:  "DEFAULT 0",
		},
		"float32": {
			Type:     "decimal(10, 3)",
			DataType: "decimal",
			Default:  "DEFAULT 0",
		},
		"time": {
			Type:     "timestamp",
			DataType: "timestamp",
			Default:  "DEFAULT CURRENT_TIMESTAMP",
		},
	}
	SpecialField = map[string]SQLField{
		`create_time`: {
			Name:    "`create_time`",
			Type:    "datetime",
			Default: "DEFAULT CURRENT_TIMESTAMP",
			Comment: "创建时间",
			NotNull: "NOT NULL",
		},
		`update_time`: {
			Name:    "`update_time`",
			Type:    "datetime",
			Default: "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
			Comment: "更新时间",
			NotNull: "NOT NULL",
		},
		`deleted_at`: {
			Name:    "`deleted_at`",
			Type:    "datetime",
			Default: "",
			Comment: "删除时间",
			NotNull: "NULL",
		},
	}
)

type GenSQL struct {
	TableName  string
	PrimaryKey string
	Engine     string
	Charset    string
	Collate    string
	Fields     []SQLField
}

type SQLField struct {
	TableName string
	After     string
	SrcName   string
	DataType  string
	Name      string
	Type      string
	Default   string
	Comment   string
	NotNull   string // 'NOT NULL' or ''
}

func (s *GenSQL) CreateTable() ([]byte, error) {
	buf := new(bytes.Buffer)
	if s.Engine == "" {
		s.Engine = "InnoDB"
	}
	if s.Charset == "" {
		s.Charset = "utf8mb4"
	}
	if s.Collate == "" {
		s.Collate = "utf8mb4_bin"
	}
	tmpl, err := template.New("GenSQL").Parse(sqlCreateTableTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *GenSQL) AddColumns() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("GenSQL").Parse(sqlCreateAddColumns)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
