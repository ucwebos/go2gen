package conf

import (
	"github.com/xbitgo/core/config"
	"github.com/xbitgo/core/di"
	"github.com/xbitgo/core/log"
)

type Config struct {
	ProjectName string `yaml:"projectName"`
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
		ConvIODir:   "{appPath}/internal/domain/converter",
		DoDir:       "{appPath}/internal/domain/repo/dbal/do",
		ConvDoDir:   "{appPath}/internal/domain/repo/dbal/converter",
		DaoDir:      "{appPath}/internal/domain/repo/dbal/dao",
		ServiceDir:  "{appPath}/internal/domain/service",
	}
)

func Init() {
	di.PrintLog = false
	cfg := config.Yaml{ConfigFile: "go2gen.yaml"}
	err := cfg.Apply(Global)
	if err != nil {
		log.Panic(err)
	}
	if Global.ProjectName == "" {
		Global.ProjectName = "miman"
	}
	if Global.Tmpl == nil {
		Global.Tmpl = defaultTmpl
	}
}
