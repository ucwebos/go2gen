package app

import (
	"fmt"
	"go2gen/conf"
	"go2gen/domian/gen"
)

type App struct {
	Pwd      string
	Name     string
	RootPath string
	Tmpl     *conf.Tmpl
}

func NewApp(rootPath string, appName string) *App {
	return &App{
		Name:     appName,
		RootPath: rootPath,
		Tmpl:     conf.Global.GetRealTmpl(rootPath, appName),
	}
}

func (a *App) parseDirs() (dirs []string) {
	dirs = []string{
		a.Tmpl.ConfDir,
		a.Tmpl.EntityDir,
		a.Tmpl.ServiceDir,
		a.Tmpl.RepoDbalDir,
	}
	dirTmpUniMap := map[string]struct{}{}
	for _, dir := range dirs {
		fmt.Println(dir)
		dirTmpUniMap[dir] = struct{}{}
	}
	items := a.ScanDir(a.Tmpl.DomainDir)
	for _, item := range items {
		fmt.Println(item)
		if _, ok := dirTmpUniMap[item]; ok {
			continue
		}
		dirTmpUniMap[item] = struct{}{}
		dirs = append(dirs, item)
	}
	return dirs
}

func (a *App) CRepo(entity string) error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.CRepo(entity)
}

func (a *App) CService(entity string) error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.CService(entity)
}

func (a *App) Tests() error {
	return nil
}

func (a *App) Dao() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	return gm.Dao()
}

func (a *App) Conv() error {
	// todo
	return nil
}
