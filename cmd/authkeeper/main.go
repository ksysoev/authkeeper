package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksysoev/authkeeper/pkg/cmd"
)

var (
	version = "dev"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Create root command with all dependencies
	rootCmd, err := cmd.NewRootCommand(ctx, version)
	if err != nil {
		return err
	}

	// Execute command
	return rootCmd.Execute()
}
