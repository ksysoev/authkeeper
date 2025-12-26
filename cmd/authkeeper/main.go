package main

import (
	"fmt"
	"os"

	"github.com/ksysoev/authkeeper/pkg/cmd"
	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/ksysoev/authkeeper/pkg/prov"
	"github.com/ksysoev/authkeeper/pkg/repo"
	"github.com/ksysoev/authkeeper/pkg/ui"
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
	// Special case: if no command provided or asking for help/version, don't prompt for password
	if len(os.Args) == 1 {
		// No command provided, show help
		vaultPath, _ := repo.GetDefaultVaultPath()
		repository := repo.NewVaultRepository(vaultPath, "dummy")
		provider := prov.NewOAuthProvider()
		service := core.NewService(repository, provider)
		cli := ui.NewCLI(service)
		app := cmd.NewApp(cli)
		rootCmd := app.BuildRootCommand(version)
		return rootCmd.Execute()
	}
	
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" || arg == "-v" || arg == "--version" || arg == "help" {
			// Build minimal app for help
			vaultPath, _ := repo.GetDefaultVaultPath()
			repository := repo.NewVaultRepository(vaultPath, "dummy")
			provider := prov.NewOAuthProvider()
			service := core.NewService(repository, provider)
			cli := ui.NewCLI(service)
			app := cmd.NewApp(cli)
			rootCmd := app.BuildRootCommand(version)
			return rootCmd.Execute()
		}
	}

	// Get vault path
	vaultPath, err := repo.GetDefaultVaultPath()
	if err != nil {
		return fmt.Errorf("failed to get vault path: %w", err)
	}

	// Check if vault exists
	tempRepo := repo.NewVaultRepository(vaultPath, "")
	vaultExists := tempRepo.Exists()
	
	// Check if user is trying to run a command that requires vault
	isAddCommand := len(os.Args) > 1 && os.Args[1] == "add"
	
	// If vault doesn't exist and it's not the 'add' command, show helper message
	if !vaultExists && !isAddCommand {
		// Build minimal app without password to show the vault-not-found message
		repository := repo.NewVaultRepository(vaultPath, "")
		provider := prov.NewOAuthProvider()
		service := core.NewService(repository, provider)
		cli := ui.NewCLI(service)
		app := cmd.NewApp(cli)
		rootCmd := app.BuildRootCommand(version)
		return rootCmd.Execute()
	}
	
	// Prompt for password only when vault exists or user is adding first client
	tempCLI := &ui.CLI{}
	password, err := tempCLI.PromptMasterPassword(!vaultExists)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Wire up dependencies (Dependency Injection)
	// Outbound adapters
	repository := repo.NewVaultRepository(vaultPath, password)
	provider := prov.NewOAuthProvider()

	// Core business logic
	service := core.NewService(repository, provider)

	// Inbound adapter
	cli := ui.NewCLI(service)

	// Application
	app := cmd.NewApp(cli)

	// Build and execute command
	rootCmd := app.BuildRootCommand(version)
	return rootCmd.Execute()
}
