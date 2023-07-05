package tpls

import (
	"bytes"
	"text/template"
)

const doTpl = `
type {{.Name}}Do struct {
{{- range .Fields}}
	{{- if eq .SType 4}}
	{{.Name}} *{{.Type}} {{.Tag}} // {{.Comment}}
	{{- else if gt .SType 0}}
	{{.Name}} string {{.Tag }} // {{.Comment}}
	{{- else}}
	{{.Name}} {{.Type}} {{.Tag}} // {{.Comment}}
	{{- end}}
{{- end}}
	DeletedAt gorm.DeletedAt ` + "`" + `db:"deleted_at" gorm:"column:deleted_at"` + "`" + ` // 软删除标识
}

func (do *{{.Name}}Do) TableName() string {
	return TableName{{.Name}}Do
}
`

type Do struct {
	Name   string
	Fields []DoField
}

type DoField struct {
	Name      string
	Type      string
	Type2     string
	SType     int
	Tag       string
	ConvSlice bool
	IsPoint   bool
	Comment   string
}

func (s *Do) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New(s.Name + "DO").Parse(doTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const convDoTpl = `

func From{{.Name}}Entity(input *entity.{{.Name}}) *do.{{.Name}}Do{
	if input == nil {
		return nil
	}
	output := &do.{{.Name}}Do{}
{{- range .Fields }}
	{{- if eq .SType 1}} 
		{{- if .IsPoint}} 
	if input.{{.Name}} != nil {
		b, _ := json.Marshal(input.{{.Name}})
		output.{{.Name}} = string(b)
	}
		{{- else}}
	b, _ := json.Marshal(input.{{.Name}})
	output.{{.Name}} = string(b)
		{{- end}}
	
	{{- else if eq .SType 2}}
	if input.{{.Name}} != nil {
		{{- if .ConvSlice}}
			output.{{.Name}} = slice_utils.Implode(input.{{.Name}},",")
		{{- else}}
		b, _ := json.Marshal(input.{{.Name}})
		output.{{.Name}} = string(b)
		{{- end}}
	}
	{{- else if eq .SType 3}}
	if input.{{.Name}} != nil {
		b, _ := json.Marshal(input.{{.Name}})
		output.{{.Name}} = string(b)
	}
	{{- else if eq .SType 4}}
		if !input.{{.Name}}.IsZero() {
			output.{{.Name}} = &input.{{.Name}}
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func To{{.Name}}Entity(input *do.{{.Name}}Do) *entity.{{.Name}}{
	if input == nil {
		return nil
	}
	output := &entity.{{.Name}}{}
{{- range .Fields }}
	{{- if eq .SType 1}} 
	if input.{{.Name}} != ""  {
		{{- if .IsPoint}} 
		t := &entity.{{ .Type2}}{}
		{{- else}}
		t := entity.{{ .Type2}}{}
		{{- end}}
		err := json.Unmarshal([]byte(input.{{.Name}}), &t)
		if err != nil {
			logx.Errorf("converter To{{$.Name}}Entity[{{.Name}}] err %v", err)
		} else {
			output.{{.Name}} = t
		}
	}
	{{- else if eq .SType 2}}
		if input.{{.Name}} != "" {
			{{- if .ConvSlice}}
				{{- if eq .Type2 "int64" }}
					output.{{.Name}} = slice_utils.ExplodeInt64(input.{{.Name}},",")
				{{- else if eq .Type2 "int" }}
					output.{{.Name}} = slice_utils.ExplodeInt(input.{{.Name}},",")
				{{- else}}
					output.{{.Name}} = slice_utils.ExplodeStr(input.{{.Name}},",")
				{{- end}}
			{{- else}}
				t := {{.Type}}{}
				err := json.Unmarshal([]byte(input.{{.Name}}), &t)
				if err != nil {
					logx.Errorf("converter To{{$.Name}}Entity[{{.Name}}] err %v", err)
				} else {
					output.{{.Name}} = t
				}
			{{- end}}
		}
	{{- else if eq .SType 3}}
		if input.{{.Name}} != "" {
			t := {{.Type}}{}
			err := json.Unmarshal([]byte(input.{{.Name}}), &t)
			if err != nil {
				logx.Errorf("converter To{{$.Name}}Entity[{{.Name}}] err %v", err)
			} else {
				output.{{.Name}} = t
			}
		}
	{{- else if eq .SType 4}}
		if input.{{.Name}} != nil {
			output.{{.Name}} = *input.{{.Name}}
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func From{{.Name}}List(input entity.{{.Name}}List) do.{{.Name}}DoList {
	if input == nil {
		return nil
	}
	output := make([]*do.{{.Name}}Do, 0, len(input))
	for _, item := range input {
		resultItem := From{{.Name}}Entity(item)
		output = append(output, resultItem)
	}
	return output
}

func To{{.Name}}List(input do.{{.Name}}DoList) entity.{{.Name}}List {
	if input == nil || len(input) == 0 {
		return nil
	}
	output := make(entity.{{.Name}}List, 0, len(input))
	for _, item := range input {
		resultItem := To{{.Name}}Entity(item)
		output = append(output, resultItem)
	}
	return output
}

`

type DoConv struct {
	Name   string
	Fields []DoField
}

func (s *DoConv) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(s.Name + "DOConv").Parse(convDoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
