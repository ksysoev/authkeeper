package cmd

import (
	"context"

	"github.com/ksysoev/authkeeper/pkg/ui"
	"github.com/spf13/cobra"
)

// App holds application dependencies
type App struct {
	cli *ui.CLI
}

// NewApp creates a new application
func NewApp(cli *ui.CLI) *App {
	return &App{
		cli: cli,
	}
}

// BuildRootCommand builds the root cobra command
func (app *App) BuildRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "authkeeper",
		Short:   "OAuth2/OIDC credential manager",
		Long:    "A beautiful CLI tool for managing OAuth2/OIDC credentials and issuing access tokens with encrypted vault storage.",
		Version: version,
	}

	rootCmd.AddCommand(app.addCommand())
	rootCmd.AddCommand(app.tokenCommand())
	rootCmd.AddCommand(app.listCommand())
	rootCmd.AddCommand(app.deleteCommand())

	return rootCmd
}

func (app *App) addCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new OIDC client to the vault",
		Long:  `Interactive command to add a new OIDC client with credentials to the encrypted vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.cli.AddClient(context.Background())
		},
	}
}

func (app *App) tokenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Issue an access token",
		Long:  `Select an OIDC client from the vault and issue an access token using client credentials flow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.cli.IssueToken(context.Background())
		},
	}
}

func (app *App) listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all OIDC clients",
		Long:  `Display all OIDC clients stored in the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.cli.ListClients(context.Background())
		},
	}
}

func (app *App) deleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete an OIDC client",
		Long:  `Select and delete an OIDC client from the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.cli.DeleteClient(context.Background())
		},
	}
}
