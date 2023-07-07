package project

import (
	"github.com/ucwebos/go2gen/conf"
	"github.com/ucwebos/go2gen/domian/app"
	"github.com/ucwebos/go2gen/domian/gen"
	"github.com/xbitgo/core/di"
	"log"
)

type Project struct {
	Pwd        string
	activeApps []*app.App
}

func NewProject(pwd string) *Project {
	p := &Project{
		Pwd:        pwd,
		activeApps: make([]*app.App, 0),
	}
	// di注册project
	di.Register("project", p)
	return p
}

func (p *Project) SetActiveApps(apps ...string) int {

	if len(apps) >= 0 {
		for _, s := range apps {
			actApp := app.NewApp(p.Pwd, s)
			p.activeApps = append(p.activeApps, actApp)
		}
	}
	return len(p.activeApps)
}

func (p *Project) RootPath() string {
	return p.Pwd
}

// Create 创建新应用
func (p *Project) Create(name string) {
	gm := gen.NewManager(conf.Global.Tmpl, name, "")
	_ = gm.App(p.Pwd, name)
}

// Generate 生成所有GO代码
func (p *Project) Generate() {
	for _, actApp := range p.activeApps {
		if err := actApp.Generate(); err != nil {
			log.Panicf("Watch app[%s] err: %v", actApp.Name, err)
		}
	}
}

// CRepo 生成
func (p *Project) CRepo(entity string) {
	for _, actApp := range p.activeApps {
		if err := actApp.CRepo(entity); err != nil {
			log.Panicf("Impl app[%s] err: %v", actApp.Name, err)
		}
	}
}

// CService 生成
func (p *Project) CService(entity string) {
	for _, actApp := range p.activeApps {
		if err := actApp.CService(entity); err != nil {
			log.Panicf("Impl app[%s] err: %v", actApp.Name, err)
		}
	}
}

// Tests 生成
func (p *Project) Tests() {
	//for _, actApp := range p.activeApps {
	//if err := actApp.Tests(); err != nil {
	//	log.Panicf("Tests app[%s] err: %v", actApp.Name, err)
	//}
	//}
}

// Dao 生成
func (p *Project) Dao() {
	for _, actApp := range p.activeApps {
		if err := actApp.Dao(); err != nil {
			log.Panicf("Dao app[%s] err: %v", actApp.Name, err)
		}
	}
}
