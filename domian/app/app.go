package app

import (
	"fmt"
	"github.com/ucwebos/go2gen/conf"
	"github.com/ucwebos/go2gen/domian/gen"
	"os"
)

type App struct {
	Pwd      string
	Name     string
	RootPath string
	Tmpl     *conf.Tmpl
}

func NewApp(rootPath string, appName string) *App {
	return &App{
		Pwd:      fmt.Sprintf("%s/%s", rootPath, appName),
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
		dirTmpUniMap[dir] = struct{}{}
	}
	items := a.ScanDir(a.Tmpl.DomainDir)
	for _, item := range items {
		if _, ok := dirTmpUniMap[item]; ok {
			continue
		}
		dirTmpUniMap[item] = struct{}{}
		dirs = append(dirs, item)
	}
	return dirs
}

func (a *App) getEntries() (entries []string) {
	fileInfos, err := os.ReadDir(a.Tmpl.EntryDir)
	if err != nil {
		return
	}
	dirList := make([]string, 0)
	for _, fi := range fileInfos {
		if fi.IsDir() {
			dirList = append(dirList, fi.Name())
		}
	}
	return dirList
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
