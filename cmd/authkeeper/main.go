package main

import (
	"os"

	"github.com/ksysoev/authkeeper/pkg/cmd"
)

var (
	version = "dev"
	name    = "authkeeper"
)

func main() {
	os.Exit(run())
}

func run() int {
	command := cmd.InitCommand(cmd.BuildInfo{
		Version: version,
		AppName: name,
	})

	if err := command.Execute(); err != nil {
		return 1
	}

	return 0
}
