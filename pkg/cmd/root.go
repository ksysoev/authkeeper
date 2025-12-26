package cmd

import (
	"context"
	"fmt"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/ksysoev/authkeeper/pkg/prov"
	"github.com/ksysoev/authkeeper/pkg/repo"
	"github.com/ksysoev/authkeeper/pkg/ui"
	"github.com/spf13/cobra"
)

// NewRootCommand creates and configures the root command with all dependencies
func NewRootCommand(ctx context.Context, version string) (*cobra.Command, error) {
	// Get vault path
	vaultPath, err := repo.GetDefaultVaultPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get vault path: %w", err)
	}

	// Check if vault exists
	tempRepo := repo.NewVaultRepository(vaultPath, "")
	vaultExists := tempRepo.Exists()

	// Prompt for password if vault exists
	var password string
	if vaultExists {
		tempCLI := &ui.CLI{}
		password, err = tempCLI.PromptMasterPassword(false)
		if err != nil {
			return nil, fmt.Errorf("failed to get password: %w", err)
		}
	}

	// Build command with dependencies
	repository := repo.NewVaultRepository(vaultPath, password)
	provider := prov.NewOAuthProvider()
	service := core.NewService(repository, provider)
	cli := ui.NewCLI(service)
	app := NewApp(cli, vaultPath)

	return app.BuildRootCommand(version), nil
}
