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
	var name, clientID, clientSecret, tokenURL, scopes string

	cmd := &cobra.Command{
		Use:   "add [flags]",
		Short: "Add a new OIDC client to the vault",
		Long:  `Add a new OIDC client with credentials to the encrypted vault. If flags are not provided, the command will prompt interactively.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cli := initCLI(arg)

			return cli.AddClient(cmd.Context(), name, clientID, clientSecret, tokenURL, scopes)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Client name")
	cmd.Flags().StringVarP(&clientID, "client-id", "c", "", "Client ID")
	cmd.Flags().StringVarP(&clientSecret, "client-secret", "s", "", "Client secret")
	cmd.Flags().StringVarP(&tokenURL, "token-url", "t", "", "Token URL")
	cmd.Flags().StringVar(&scopes, "scopes", "", "Scopes (space-separated)")

	return cmd
}

// TokenCommand creates a new cobra.Command to issue an access token.
// It returns a pointer to a cobra.Command which can be executed to issue a token.
func TokenCommand(arg *args) *cobra.Command {
	var clientName string

	cmd := &cobra.Command{
		Use:   "token [client-name]",
		Short: "Issue an access token",
		Long:  `Issue an access token for an OIDC client using client credentials flow. If client name is not provided, you will be prompted to select from available clients.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := initCLI(arg)

			if len(args) > 0 && clientName == "" {
				clientName = args[0]
			}

			return cli.IssueToken(cmd.Context(), clientName)
		},
	}

	cmd.Flags().StringVarP(&clientName, "client", "c", "", "Client name")

	return cmd
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
	var clientName string
	var force bool

	cmd := &cobra.Command{
		Use:   "delete [client-name]",
		Short: "Delete an OIDC client",
		Long:  `Delete an OIDC client from the vault. If client name is not provided, you will be prompted to select from available clients.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := initCLI(arg)

			if len(args) > 0 && clientName == "" {
				clientName = args[0]
			}

			return cli.DeleteClient(cmd.Context(), clientName, force)
		},
	}

	cmd.Flags().StringVarP(&clientName, "client", "c", "", "Client name")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
