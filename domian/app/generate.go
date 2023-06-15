package app

import (
	"fmt"
	"go2gen/domian/gen"
	"go2gen/domian/parser"
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
	fmt.Println(dirs)
	for _, dir := range dirs {
		pkg := &Pkg{
			Dir: dir,
		}
		if dir == a.Tmpl.EntityDir {
			pkg.IsEntity = true
		}
		fmt.Println(dir)
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
	// entity 生成
	for _, xst := range entityPkg.Parser.StructList {
		err := gm.Do(xst)
		if err != nil {
			log.Printf("gen do err: %v \n", err)
			return err
		}
	}
	gm.EntityTypeDef()
	// do 的TypeDef
	gm.DoTypeDef()

	if err := a.Dao(); err != nil {
		log.Printf("Dao app[%s] err: %v", a.Name, err)
		return err
	}
	if err := a.Conv(); err != nil {
		log.Printf("Conv app[%s] err: %v", a.Name, err)
		return err
	}

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
