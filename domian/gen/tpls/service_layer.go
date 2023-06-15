package tpls

import (
	"bytes"
	"text/template"
)

const serviceTpl = `
package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xbitgo/components/dtx"
	"github.com/xbitgo/components/filterx"
	
	"{{.ProjectName}}/apps/{{.AppName}}/domain/entity"
	"{{.ProjectName}}/apps/{{.AppName}}/domain/repo"
)

// {{.ServiceName}} @DI
type {{.ServiceName}} struct {
	{{- range .EntityList }}
	{{.EntityName}}Repo repo.{{.EntityName}}Repo ` + "`" + `di:"repo_impl.{{.EntityName}}RepoImpl"` + "`" + `
	{{- end}}
}

func New{{.ServiceName}}() *{{.ServiceName}} {
	return &{{.ServiceName}}{}
}


{{- range .EntityList }}
func (s *{{$.ServiceName}}) Create{{.EntityName}}(ctx context.Context, {{.VarName}} *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	_{{.VarName}}, err := s.{{.EntityName}}Repo.Create(ctx, {{.VarName}})
	if err != nil {
		return nil, err
	}
	return _{{.VarName}}, nil
}

func (s *{{$.ServiceName}}) Query{{.EntityName}}(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	list, count, err := s.{{.EntityName}}Repo.Query(ctx, query, pg)
	if err != nil {
		return nil, 0, err
	}
	return list, count, err
}

func (s *{{$.ServiceName}}) Get{{.EntityName}}(ctx context.Context, _{{.VarName}}Id int64) (*entity.{{.EntityName}}, error) {
	_{{.VarName}}, err := s.{{.EntityName}}Repo.Get(ctx, _{{.VarName}}Id)
	if err != nil {
		return nil, err
	}
	return _{{.VarName}}, nil
}

func (s *{{$.ServiceName}}) Set{{.EntityName}}(ctx context.Context, {{.VarName}}Id int64, updateMap map[string]interface{}) (*entity.{{.EntityName}}, error) {
	setItems := dtx.SetItemList{}
	for k, v := range updateMap {
		setItems = append(setItems, &dtx.SetItem{
			Field:    k,
			Operator: dtx.SET,
			Value:    v,
		})
	}
	err := s.{{.EntityName}}Repo.UpdateById(ctx, setItems, {{.VarName}}Id)
	if err != nil {
		return nil, err
	}
	return s.{{.EntityName}}Repo.Get(ctx, {{.VarName}}Id)
}

func (s *{{$.ServiceName}}) Delete{{.EntityName}}(ctx context.Context, {{.VarName}}Id int64) error {
	return s.{{.EntityName}}Repo.DeleteById(ctx, {{.VarName}}Id)
}
{{- end}}
`

type ServiceLayer struct {
	ProjectName string
	AppName     string
	ServiceName string
	EntityList  []ServiceLayerItem
}

type ServiceLayerItem struct {
	EntityName string
	VarName    string
}

func (s *ServiceLayer) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("ServiceLayer").Parse(serviceTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
