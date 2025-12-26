package main

import (
	"os"

	"github.com/ksysoev/authkeeper/pkg/cmd"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	rootCmd, err := cmd.InitCommands(version)
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}
