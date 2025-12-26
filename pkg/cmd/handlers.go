package cmd

import (
	"context"
	"fmt"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/ksysoev/authkeeper/pkg/prov"
	"github.com/ksysoev/authkeeper/pkg/repo"
	"github.com/ksysoev/authkeeper/pkg/ui"
)

// handleAddClient handles the add command with vault initialization if needed
func (app *App) handleAddClient(ctx context.Context) error {
	// If vault is not initialized, initialize it first
	if !app.cli.IsVaultInitialized() {
		if err := app.initializeVault(); err != nil {
			return fmt.Errorf("failed to initialize vault: %w", err)
		}
	}

	return app.cli.AddClient(ctx)
}

// initializeVault prompts for password and reinitializes the CLI with the new vault
func (app *App) initializeVault() error {
	tempCLI := &ui.CLI{}
	password, err := tempCLI.PromptMasterPassword(true)
	if err != nil {
		return err
	}

	// Recreate repository with password
	repository := repo.NewVaultRepository(app.vaultPath, password)
	provider := prov.NewOAuthProvider()
	service := core.NewService(repository, provider)
	
	// Replace CLI service with initialized one
	*app.cli = *ui.NewCLI(service)

	return nil
}

// vaultNotInitializedError returns a helpful error message when vault is not initialized
func (app *App) vaultNotInitializedError() error {
	return fmt.Errorf("vault not initialized. Please run 'authkeeper add' to create your first client")
}
