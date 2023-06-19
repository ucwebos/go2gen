package tpls

import (
	"bytes"
	"text/template"
)

const ioTpl = `
type {{.Name}} {
{{- range .Fields}}
	{{- if eq .SType 4}}
	{{.Name}} {{.Type}} {{.Tag}} // {{.Comment}}
	{{- else if gt .SType 0}}
	{{.Name}} string {{.Tag }} // {{.Comment}}
	{{- else}}
	{{.Name}} {{.Type}} {{.Tag}} // {{.Comment}}
	{{- end}}
{{- end}}
}
`

type IO struct {
	Name   string
	Fields []IoField
}

type IoField struct {
	Name    string
	Type    string
	Type2   string
	SType   int
	Tag     string
	Hidden  bool
	Comment string
}

func (s *IO) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New(s.Name + "IO").Parse(ioTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const convIoTpl = `
func From{{.Name}}Entity(input *entity.{{.Name}}) *types.{{.Name}}{
	if input == nil {
		return nil
	}
	output := &types.{{.Name}}{}
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
			output.{{.Name}} = tool_time.TimeToDateTimeString(input.{{.Name}})
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func To{{.Name}}Entity(input *types.{{.Name}}) *entity.{{.Name}}{
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
			log.Errorf("converter To{{$.Name}}Entity[{{.Name}}] err %v", err)
		} else {
			output.{{.Name}} = t
		}
	}
	{{- else if eq .SType 2}}
		if input.{{.Name}} != "" {
			output.{{.Name}} = input.{{.Name}}
		}
	{{- else if eq .SType 3}}
		if input.{{.Name}} != "" {
			//t := {{.Type}}{}
			output.{{.Name}} = To{{.Name}}Entity(input.{{.Name}})
		}
	{{- else if eq .SType 4}}
		if ts := tool_time.ParseDateTime(input.{{.Name}}); !ts.IsZero() {
			output.{{.Name}} = ts
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func From{{.Name}}List(input entity.{{.Name}}List) []*types.{{.Name}} {
	if input == nil {
		return nil
	}
	output := make([]*types.{{.Name}}, 0, len(input))
	for _, item := range input {
		resultItem := From{{.Name}}Entity(item)
		output = append(output, resultItem)
	}
	return output
}

func To{{.Name}}List(input []*types.{{.Name}}) entity.{{.Name}}List {
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

type IoConv struct {
	SrcPath string
	Name    string
	Package string
	Imports []string
	Fields  []IoField
}

func (s *IoConv) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New(s.Name + "IOConv").Parse(convIoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
