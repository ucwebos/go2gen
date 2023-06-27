package tpls

import (
	"bytes"
	"go2gen/domian/parser"
	"text/template"
)

const HandlerFuncInitTpl = `package handler

import (
	"github.com/gin-gonic/gin"

	"{{.ProjectName}}/{{.AppName}}/internal/entry/{{.Entry}}/types"
)
{{range .FuncList}}
func {{.Key}}(ctx *gin.Context{{if ne .Request nil}}, req *types.{{.Request.Name}}{{end}}) ({{if ne .Response nil}}*types.{{.Response.Name}}, {{end}}error) {
	{{- if ne .Response nil}}
	var (
		resp = &types.{{.Response.Name}}{}
	)
	{{- end}}
	// todo ...
	return {{if ne .Response nil}}resp, {{end}}nil
}
{{end}}
`

type HandlerFunc struct {
	ProjectName string
	AppName     string
	Entry       string
	FuncList    []*parser.EntryModuleFunc
}

const HandlerFuncAppendTpl = `
{{range .FuncList}}
func {{.Key}}(ctx *gin.Context, req *types.{{.Request.Name}}) (resp *types.{{.Response.Name}}, err error) {
	// todo ...
	return resp, nil
}
{{end}}
`

type HandlerFuncAppend struct {
	Body     []byte
	FuncList []*parser.EntryModuleFunc
}

func (s *HandlerFuncAppend) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerFuncAppend").Parse(HandlerFuncAppendTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}

func (s *HandlerFunc) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerFunc").Parse(HandlerFuncInitTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
