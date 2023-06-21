package tpls

import (
	"bytes"
	"text/template"
)

const repoTpl = `
package repo

import (
	"context"

	"{{.ProjectName}}/{{.AppName}}/internal/domain/entity"
	"{{.ProjectName}}/{{.AppName}}/internal/domain/repo/dbal"

	"{{.ProjectName}}/common/tools/filterx"
)

// {{.EntityName}}Repo . @GI
type {{.EntityName}}Repo struct {
	DBAL *dbal.{{.EntityName}}RepoDBAL
}

func New{{.EntityName}}Repo() *{{.EntityName}}Repo {
	return &{{.EntityName}}Repo{
		DBAL: dbal.New{{.EntityName}}RepoDBAL(),
	}
}

func (r *{{.EntityName}}Repo) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	return r.DBAL.Query(ctx,query,pg)
}

func (r *{{.EntityName}}Repo) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	return r.DBAL.Create(ctx,input)
}


{{- if .HasID}}
func (r *{{.EntityName}}Repo) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	return r.DBAL.GetByID(ctx,id)
}

func (r *{{.EntityName}}Repo) UpdateById(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.DBAL.UpdateById(ctx,id,updates)
}

func (r *{{.EntityName}}Repo) DeleteById(ctx context.Context, id int64) error {
	return r.DBAL.DeleteById(ctx,id)
}
{{- end}}

`

const RepoDBALTpl = `package dbal

import (
	"context"

	"gorm.io/gorm"

	"{{.ProjectName}}/{{.AppName}}/internal/config"
	"{{.ProjectName}}/{{.AppName}}/internal/domain/entity"
	"{{.ProjectName}}/{{.AppName}}/internal/domain/repo/dbal/converter"
	"{{.ProjectName}}/{{.AppName}}/internal/domain/repo/dbal/dao"
	"{{.ProjectName}}/{{.AppName}}/internal/domain/repo/dbal/do"

	"{{.ProjectName}}/common/lib/db"
	"{{.ProjectName}}/common/tools/filterx"
)

// {{.EntityName}}RepoDBAL .
type {{.EntityName}}RepoDBAL struct {
	DB  *db.DBWrapper 
	Dao *dao.{{.EntityName}}Dao
}

func New{{.EntityName}}RepoDBAL() *{{.EntityName}}RepoDBAL {
	return &{{.EntityName}}RepoDBAL{
		DB:  config.GetDB(),
		Dao: dao.New{{.EntityName}}Dao(),
	}
}

func (impl *{{.EntityName}}RepoDBAL) NewReadSession(ctx context.Context) *gorm.DB {
	return impl.DB.NewSession(ctx)
}

func (impl *{{.EntityName}}RepoDBAL) NewCreateSession(ctx context.Context) *gorm.DB {
	return impl.DB.NewSession(ctx)
}

func (impl *{{.EntityName}}RepoDBAL) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	session := impl.NewReadSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return nil, 0, err
	}
	session, noCount := filterx.PageGormOption(session, pg)
	var (
		doList do.{{.EntityName}}DoList
		count  int
	)
	if noCount {
		doList, err = impl.Dao.FindAll(session)
	} else {
		doList, count, err = impl.Dao.FindPage(session)
	}
	if err != nil {
		return nil, 0, err
	}
	return converter.To{{.EntityName}}List(doList), count, nil
}

func (impl *{{.EntityName}}RepoDBAL) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	session := impl.NewCreateSession(ctx)
	_do := converter.From{{.EntityName}}Entity(input)
	err := impl.Dao.Create(session, _do)
	if err != nil {
		return nil, err
	}
	output := converter.To{{.EntityName}}Entity(_do)
	return output, err
}

{{- if .HasID}}
func (impl *{{.EntityName}}RepoDBAL) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	session := impl.NewReadSession(ctx)
	session = session.Where("id = ?",id)
	_do, err := impl.Dao.Get(session)
	if err != nil {
		return nil, err
	}
	return converter.To{{.EntityName}}Entity(_do), nil
}

func (impl *{{.EntityName}}RepoDBAL) UpdateById(ctx context.Context, id int64, updates map[string]interface{}) error {
	session := impl.NewReadSession(ctx)
	session = session.Where("id = ?",id)
	err := impl.Dao.Update(session, updates)
	if err != nil {
		return err
	}
	return err
}

func (impl *{{.EntityName}}RepoDBAL) DeleteById(ctx context.Context, id int64) error {
	session := impl.NewReadSession(ctx)
	session = session.Where("id = ?",id)
	err := impl.Dao.Delete(session)
	if err != nil {
		return err
	}
	return err
}
{{- end}}
`

type Repo struct {
	ProjectName string
	AppName     string
	EntityName  string
	HasID       bool
}

func (s *Repo) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl").Parse(repoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Repo) ExecuteImpl() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl.impl").Parse(RepoDBALTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
