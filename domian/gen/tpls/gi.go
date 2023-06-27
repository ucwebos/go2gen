package tpls

import (
	"bytes"
	"text/template"
)

const GITpl = `// Code generated by go2gen. DO NOT EDIT.
package {{.Pkg}}

import (
    "log"

    "github.com/xbitgo/core/di"
)

var (
	{{- range .List}}
		{{.NameVal}}Instance *{{.Name}}
	{{- end}}
)


{{- range .List}}
func {{.Name}}Instance() *{{.Name}} {
	if {{.NameVal}}Instance == nil {
		{{- if ne .NewReturnsLen 0}}
			{{- if ge .NewReturnsLen 2 }}
			 _{{.NameVal}}, err := New{{.Name}}()
			 if err != nil {
				log.Panicf("init GI {{.Pkg}}.{{.Name}}] err: %v", err)
			 }
			{{- else}}
				_{{.NameVal}} := New{{.Name}}()
			{{- end}}
		{{- else}}
			_{{.NameVal}} := &{{.Name}}{}
		{{- end}}
		{{.NameVal}}Instance = _{{.NameVal}}
	}
	return {{.NameVal}}Instance
}
{{- end}}

`

type DI struct {
	Pkg  string // 包名
	List []DItem
}

func NewGI(pkg string) *DI {
	return &DI{
		Pkg:  pkg,
		List: []DItem{},
	}
}

func (d *DI) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("GI").Parse(GITpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
