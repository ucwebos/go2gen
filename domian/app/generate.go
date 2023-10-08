package app

import (
	"github.com/ucwebos/go2gen/conf"
	"github.com/ucwebos/go2gen/domian/gen"
	"github.com/ucwebos/go2gen/domian/parser"
	"log"
)

type Pkg struct {
	Dir      string
	IsEntity bool
	Parser   *parser.IParser
}

func (a *App) Generate() error {
	var (
		dirs    = a.parseDirs()
		pkgList = map[string]*Pkg{}
	)
	for _, dir := range dirs {
		pkg := &Pkg{
			Dir: dir,
		}
		if dir == a.Tmpl.EntityDir {
			pkg.IsEntity = true
		}
		pr, err := parser.Scan(dir, parser.ParseTypeWatch)
		if err != nil {
			log.Printf("generate parser pkg[%s] err: %v \n", dir, err)
			return err
		}
		pkg.Parser = pr
		pkgList[dir] = pkg
	}

	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	entityPkg := pkgList[a.Tmpl.EntityDir]

	gm.DoList(entityPkg.Parser.StructList)

	for _, entry := range a.getEntries() {
		gm.IOEntries(entityPkg.Parser.StructList, entityPkg.Parser.AliasList, entry)
	}

	gm.EntityTypeDef()
	// do 的TypeDef
	gm.DoTypeDef()

	// c.repo
	for _, xst := range entityPkg.Parser.StructList {
		a.CRepo(xst.Name)
	}

	// dao
	if err := a.Dao(); err != nil {
		log.Printf("Dao app[%s] err: %v", a.Name, err)
		return err
	}

	// conv
	if err := a.Conv(); err != nil {
		log.Printf("Conv app[%s] err: %v", a.Name, err)
		return err
	}

	// handler
	a.BatchHandlerAndDoc()

	// GI
	if err := a.BatchGI(); err != nil {
		return err
	}

	// sql
	//dsn := conf.Global.DB
	//if dsn != "" {
	//	gm.Do2Sql(dsn)
	//}

	return nil
}

func (a *App) BatchHandlerAndDoc() error {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	for _, entry := range a.getEntries() {
		gm.HandlerAndDoc(entry)
	}
	return nil
}

func (a *App) BatchGI() error {
	var (
		dirs    = a.parseDirs()
		pkgList = map[string]*Pkg{}
	)
	for _, dir := range dirs {
		pkg := &Pkg{
			Dir: dir,
		}
		if dir == a.Tmpl.EntityDir {
			pkg.IsEntity = true
		}
		pr, err := parser.Scan(dir, parser.ParseTypeWatch)
		if err != nil {
			log.Printf("generate parser pkg[%s] err: %v \n", dir, err)
			return err
		}
		pkg.Parser = pr
		pkgList[dir] = pkg
	}
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	// 其他的 DI生成
	for s, pkg := range pkgList {
		if s == a.Tmpl.EntityDir {
			continue
		}
		if s == a.Tmpl.ConfDir {
			continue
		}
		err := gm.GI(pkg.Parser)
		if err != nil {
			log.Printf("gen di err: %v \n", err)
			return err
		}
	}
	return nil
}

func (a *App) Do2SQL() {
	gm := gen.NewManager(a.Tmpl, a.Name, a.Pwd)
	dsn := conf.Global.DB
	if dsn != "" {
		gm.Do2Sql(dsn)
	}
}
