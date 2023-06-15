package main

import (
	"log"

	"github.com/spf13/cobra"

	"go2gen/conf"
	"go2gen/handler"
)

var rootCmd = &cobra.Command{
	Use:     "go2gen",
	Short:   "go2gen: An toolkit for golang code generate.",
	Long:    "go2gen: An toolkit for golang code generate..",
	Version: "0.0.1",
}

func init() {
	conf.Init()
	rootCmd.AddCommand(handler.CmdList()...)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
