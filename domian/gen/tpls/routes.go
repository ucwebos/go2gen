package tpls

import (
	"bytes"
	"go2gen/domian/parser"
	"text/template"
)

const RoutesTpl = `package api

import (
	"github.com/gin-gonic/gin"

	"{{.ProjectName}}/{{.AppName}}/internal/entry/{{.Entry}}/handler"
	"{{.ProjectName}}/{{.AppName}}/internal/entry/{{.Entry}}/types"
	"{{.ProjectName}}/{{.AppName}}/internal/middleware"
	"{{.ProjectName}}/common"
)

func generated(r *gin.Engine) {
	{{range $it := .ModuleList}}
	// {{$it.Name}}
	r.Use({{range $it.Middleware}}middleware.{{.}}(){{end}})
	{
		{{- range $v := $it.FuncList}}
		// {{$v.Name}}
		r.POST("/{{$.Entry}}/{{$it.Key}}/{{$v.KeyLi}}", func(ctx *gin.Context) {
			_raw, err := commonVerify(ctx)
			if err != nil {
				common.JSONError(ctx, err.(common.ErrCode))
				return
			}
			{{- if ne $v.Request nil}}
			var req = &types.{{$v.Request.Name}}{}
			if err := common.BindBody(_raw, &req); err != nil {
				common.JSONError(ctx, common.ErrParams)
				return
			}
			{{- if ne $v.Response nil}}
			res, err := handler.{{$v.Key}}(ctx, req)
			common.JSON(ctx, res, err)
			{{- else}}
			err := handler.{{$v.Key}}(ctx, req)
			common.JSON(ctx, nil, err)
			{{- end}}
			{{- else}}
			res, err := handler.{{$v.Key}}(ctx)
			common.JSON(ctx, res, err)
			{{- end}}
		})

		{{- end}}
	}
	{{end}}
}
`

type Routes struct {
	ProjectName string
	AppName     string
	Entry       string
	ModuleList  []parser.EntryModule
}

func (s *Routes) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("ROUTES").Parse(RoutesTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
