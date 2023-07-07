package handler

import (
	"github.com/spf13/cobra"
	"github.com/ucwebos/go2gen/domian/project"
	"log"
	"os"
)

var pwd, _ = os.Getwd()

func CmdList() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "generate",
			Short: "生成所有go代码 依次 do > dao > c.repo > conv > GI",
			Long:  "生成所有go代码 依次 do > dao > c.repo > conv > GI",
			Run:   generate,
		},
		{
			Use:   "tests",
			Short: "生成接口单元测试用例",
			Long:  "生成接口单元测试用例; 参数 {app}; app为应用名称 必须",
			Run:   tests,
		},
	}
}

func generate(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num == 0 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Generate()
}

func tests(cmd *cobra.Command, args []string) {
	p := project.NewProject(pwd)
	num := p.SetActiveApps(args...)
	if num != 1 {
		log.Fatalf("请输入正确的应用名! ")
	}
	p.Tests()
}
