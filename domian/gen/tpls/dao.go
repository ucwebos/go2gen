package tpls

import (
	"bytes"
	"text/template"
)

const daoTpl = `
type {{.DaoName}} struct {
}

func New{{.DaoName}}() *{{.DaoName}} {
  return &{{.DaoName}}{}
}

{{if .PkName}}func (dao *{{.DaoName}}) GetById(db *gorm.DB, id {{.PkType}}) (*do.{{.EntityName}}, error) {
	result := &do.{{.EntityName}}{}
	err := db.Table({{.TableName}}).Where("{{.PkCol}} = ?", id).First(result).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} GetById failed")
	}
	return result, nil
} 
{{end}}

{{if .PkName}}func (dao *{{.DaoName}}) GetByIdList(db *gorm.DB, idList []{{.PkType}}) (do.{{.EntityListName}}, error) {
	result := make([]*do.{{.EntityName}}, 0)
	if err := db.Table({{.TableName}}).Where("{{.PkCol}} in (?)", idList).Find(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} GetByIdList failed")
	}
	return result, nil
} 
{{end}}

func (dao *{{.DaoName}}) Create(db *gorm.DB, data *do.{{.EntityName}}) error {
	err := db.Table({{.TableName}}).Create(data).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Create failed")
	}
	return nil
}

func (dao *{{.DaoName}}) Save(db *gorm.DB, data *do.{{.EntityName}}) error {
	err := db.Table({{.TableName}}).Save(data).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Save failed")
	}
	return nil
}


func (dao *{{.DaoName}}) CreateBatch(db *gorm.DB, data do.{{.EntityListName}}) error {
	err := db.Table({{.TableName}}).CreateInBatches(data, len(data)).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} CreateBatch failed")
	}
	return nil
}

func (dao *{{.DaoName}}) Update(db *gorm.DB,updates map[string]any) error {
	var err error = nil
	db = db.Table({{.TableName}})
	if db.Updates(updates).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} Update failed")
	}
	return nil
}

func (dao *{{.DaoName}}) Delete(db *gorm.DB) error {
	var err error = nil
	if db.Delete(&do.{{.EntityName}}{}).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} DeleteById failed")
	}
	return nil
}

{{if .PkName}}func (dao *{{.DaoName}}) UpdateById(db *gorm.DB,updates map[string]any,id int64) error {
	var err error = nil
	db = db.Table({{.TableName}}).Where("{{.PkCol}} = ?", id)
	if db.Updates(updates).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} UpdateById failed")
	}
	return nil
}

func (dao *{{.DaoName}}) UpdateByIdList(db *gorm.DB,updates map[string]any,idList []int64) error {
	var err error = nil
	db = db.Table({{.TableName}}).Where("{{.PkCol}} in ?", idList)
	if db.Updates(updates).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} UpdateByIdList failed")
	}
	return nil
}
func (dao *{{.DaoName}}) DeleteById(db *gorm.DB,id int64) error {
	var err error = nil
	if db.Delete(&do.{{.EntityName}}{},"{{.PkCol}} = ?", id).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} DeleteById failed")
	}
	return nil
}

func (dao *{{.DaoName}}) DeleteByIdList(db *gorm.DB,idList []int64) error {
	var err error = nil
	if db.Delete(&do.{{.EntityName}}{},"{{.PkCol}} in ?", idList).Error != nil {
		return errors.Wrapf(err, "{{.DaoName}} DeleteByIdList failed")
	}
	return nil
}
{{end}}

func (dao *{{.DaoName}}) FindPage(db *gorm.DB) (do.{{.EntityListName}}, int, error) {
	db.Table({{.TableName}})
	result := make([]*do.{{.EntityName}}, 0)
	err := db.Find(&result).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.DaoName}} FindPage failed 数据库错误")
	}
	delete(db.Statement.Clauses, "LIMIT")
	var count int64
	err = db.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.DaoName}} FindPage failed 数据库错误")
	}
	return result, int(count), nil
}

func (dao *{{.DaoName}}) FindAll(db *gorm.DB) (do.{{.EntityListName}}, error) {
	result := make([]*do.{{.EntityName}}, 0)
	err := db.Find(&result).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} FindAll failed 数据库错误")
	}
	return result, nil
}

func (dao *{{.DaoName}}) Get(db *gorm.DB) (*do.{{.EntityName}}, error) {
    db.Table({{.TableName}})
	result := &do.{{.EntityName}}{}
	err := db.First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("记录获取失败")
	}
	return result, nil
}`

type Dao struct {
	EntityName     string
	DaoName        string
	EntityListName string
	TableName      string
	PkName         string
	PkType         string
	PkCol          string
}

func (s *Dao) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Dao").Parse(daoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
