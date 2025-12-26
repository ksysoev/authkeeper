package cmd

import (
	"fmt"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/ksysoev/authkeeper/pkg/prov"
	"github.com/ksysoev/authkeeper/pkg/repo"
	"github.com/ksysoev/authkeeper/pkg/ui"
	"github.com/spf13/cobra"
)

type args struct {
	version   string
	vaultPath string
}

// InitCommands initializes and returns the root command for the AuthKeeper service.
// It sets up the command structure and adds subcommands, including setting up persistent flags.
// It returns a pointer to a cobra.Command which represents the root command.
func InitCommands(version string) (*cobra.Command, error) {
	args := &args{
		version: version,
	}

	vaultPath, err := getDefaultVaultPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get vault path: %w", err)
	}

	args.vaultPath = vaultPath

	cmd := &cobra.Command{
		Use:     "authkeeper",
		Short:   "OAuth2/OIDC credential manager",
		Long:    "A beautiful CLI tool for managing OAuth2/OIDC credentials and issuing access tokens with encrypted vault storage.",
		Version: version,
	}

	cmd.PersistentFlags().StringVarP(&args.vaultPath, "vault", "v", args.vaultPath, "Path to the encrypted vault file")

	cmd.AddCommand(AddCommand(args))
	cmd.AddCommand(TokenCommand(args))
	cmd.AddCommand(ListCommand(args))
	cmd.AddCommand(DeleteCommand(args))

	return cmd, nil
}

// initCLI initializes the CLI with the vault repository and provider
func initCLI(arg *args) *ui.CLI {
	repository := repo.NewVaultRepository(arg.vaultPath)
	provider := prov.NewOAuthProvider()
	service := core.NewService(repository, provider)

	return ui.NewCLI(service)
}

// AddCommand creates a new cobra.Command to add a new OIDC client to the vault.
// It returns a pointer to a cobra.Command which can be executed to add a client.
func AddCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new OIDC client to the vault",
		Long:  `Interactive command to add a new OIDC client with credentials to the encrypted vault.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cli := initCLI(arg)

			return cli.AddClient(cmd.Context())
		},
	}
}

// TokenCommand creates a new cobra.Command to issue an access token.
// It returns a pointer to a cobra.Command which can be executed to issue a token.
func TokenCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Issue an access token",
		Long:  `Select an OIDC client from the vault and issue an access token using client credentials flow.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cli := initCLI(arg)

			return cli.IssueToken(cmd.Context())
		},
	}
}

// ListCommand creates a new cobra.Command to list all OIDC clients.
// It returns a pointer to a cobra.Command which can be executed to list clients.
func ListCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all OIDC clients",
		Long:  `Display all OIDC clients stored in the vault.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cli := initCLI(arg)

			return cli.ListClients(cmd.Context())
		},
	}
}

// DeleteCommand creates a new cobra.Command to delete an OIDC client.
// It returns a pointer to a cobra.Command which can be executed to delete a client.
func DeleteCommand(arg *args) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete an OIDC client",
		Long:  `Select and delete an OIDC client from the vault.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cli := initCLI(arg)

			return cli.DeleteClient(cmd.Context())
		},
	}
}
