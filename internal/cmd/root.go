package cmd

import (
	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Version string
	AppName string
}

func InitCommand(info BuildInfo) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     info.AppName,
		Short:   "AuthKeeper - Secure secret manager and OAuth2/OIDC client",
		Long:    `A beautiful CLI tool for managing OAuth2/OIDC credentials and issuing access tokens with encrypted vault storage.`,
		Version: info.Version,
	}

	rootCmd.AddCommand(newAddCommand())
	rootCmd.AddCommand(newTokenCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newDeleteCommand())

	return rootCmd
}
