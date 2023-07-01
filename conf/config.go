package conf

import (
	"github.com/xbitgo/core/config"
	"github.com/xbitgo/core/di"
	"github.com/xbitgo/core/log"
)

type Config struct {
	ProjectName string `yaml:"projectName"`
	DB          string `yaml:"db"`
	Tmpl        *Tmpl  `yaml:"tmpl"`
}

var (
	Global      = &Config{}
	defaultTmpl = &Tmpl{
		ConfDir:     "{appPath}/internal/config",
		DomainDir:   "{appPath}/internal/domain",
		EntityDir:   "{appPath}/internal/domain/entity",
		RepoDir:     "{appPath}/internal/domain/repo",
		RepoDbalDir: "{appPath}/internal/domain/repo/dbal",
		DoDir:       "{appPath}/internal/domain/repo/dbal/do",
		ConvDoDir:   "{appPath}/internal/domain/repo/dbal/converter",
		DaoDir:      "{appPath}/internal/domain/repo/dbal/dao",
		SQLDir:      "{appPath}/internal/domain/repo/dbal/sql",
		ServiceDir:  "{appPath}/internal/domain/service",
		EntryDir:    "{appPath}/internal/entry",
	}
)

func Init() {
	di.PrintLog = false
	cfg := config.Yaml{ConfigFile: ".go2gen.yaml"}
	err := cfg.Apply(Global)
	if err != nil {
		log.Panic(err)
	}
	if Global.ProjectName == "" {
		Global.ProjectName = "projectName"
	}
	if Global.Tmpl == nil {
		Global.Tmpl = defaultTmpl
	}
}
