package conf

import (
	"fmt"
	"strings"
)

type Tmpl struct {
	ConfDir     string `json:"conf_dir" yaml:"conf_dir"`
	DomainDir   string `json:"domain_dir" yaml:"domain_dir"`
	EntityDir   string `json:"entity_dir" yaml:"entity_dir"`
	DoDir       string `json:"do_dir" yaml:"do_dir"`
	ConvDoDir   string `json:"conv_do_dir" yaml:"conv_do_dir"`
	DaoDir      string `json:"dao_dir" yaml:"dao_dir"`
	ConvIODir   string `json:"conv_io_dir" yaml:"conv_io_dir"`
	RepoDir     string `json:"repo_dir" yaml:"repo_dir"`
	RepoDbalDir string `json:"repo_dbal_dir" yaml:"repo_dbal_dir"`
	ServiceDir  string `json:"service_dir" yaml:"service_dir"`
}

func (c *Config) GetRealTmpl(rootPath, appName string) *Tmpl {
	return &Tmpl{
		ConfDir:     replacePath(c.Tmpl.ConfDir, rootPath, appName),
		DomainDir:   replacePath(c.Tmpl.DomainDir, rootPath, appName),
		EntityDir:   replacePath(c.Tmpl.EntityDir, rootPath, appName),
		DoDir:       replacePath(c.Tmpl.DoDir, rootPath, appName),
		ConvDoDir:   replacePath(c.Tmpl.ConvDoDir, rootPath, appName),
		DaoDir:      replacePath(c.Tmpl.DaoDir, rootPath, appName),
		ConvIODir:   replacePath(c.Tmpl.ConvIODir, rootPath, appName),
		RepoDir:     replacePath(c.Tmpl.RepoDir, rootPath, appName),
		RepoDbalDir: replacePath(c.Tmpl.RepoDbalDir, rootPath, appName),
		ServiceDir:  replacePath(c.Tmpl.ServiceDir, rootPath, appName),
	}
}

func replacePath(str string, rootPath, appName string) string {
	appPath := fmt.Sprintf("%s/%s", rootPath, appName)
	str = strings.Replace(str, "{appPath}", appPath, -1)
	return str
}
