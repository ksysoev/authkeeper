package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/authkeeper/internal/tui"
	"github.com/ksysoev/authkeeper/internal/vault"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new OIDC client to the vault",
		Long:  `Interactive command to add a new OIDC client with credentials to the encrypted vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := vault.GetDefaultVaultPath()
			if err != nil {
				return fmt.Errorf("failed to get vault path: %w", err)
			}

			v := vault.New(vaultPath)

			model := tui.NewAddClientModel(v)
			p := tea.NewProgram(model, tea.WithAltScreen())

			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to run TUI: %w", err)
			}

			if m, ok := finalModel.(tui.AddClientModel); ok {
				if m.Err != nil {
					return m.Err
				}
			}

			return nil
		},
	}
}

func newTokenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Issue an access token",
		Long:  `Select an OIDC client from the vault and issue an access token using client credentials flow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := vault.GetDefaultVaultPath()
			if err != nil {
				return fmt.Errorf("failed to get vault path: %w", err)
			}

			v := vault.New(vaultPath)

			model := tui.NewTokenModel(v)
			p := tea.NewProgram(model, tea.WithAltScreen())

			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to run TUI: %w", err)
			}

			if m, ok := finalModel.(tui.TokenModel); ok {
				if m.Err != nil {
					return m.Err
				}
			}

			return nil
		},
	}
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all OIDC clients",
		Long:  `Display all OIDC clients stored in the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := vault.GetDefaultVaultPath()
			if err != nil {
				return fmt.Errorf("failed to get vault path: %w", err)
			}

			v := vault.New(vaultPath)

			model := tui.NewListModel(v)
			p := tea.NewProgram(model, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				return fmt.Errorf("failed to run TUI: %w", err)
			}

			return nil
		},
	}
}

func newDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete an OIDC client",
		Long:  `Select and delete an OIDC client from the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := vault.GetDefaultVaultPath()
			if err != nil {
				return fmt.Errorf("failed to get vault path: %w", err)
			}

			v := vault.New(vaultPath)

			model := tui.NewDeleteModel(v)
			p := tea.NewProgram(model, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			return nil
		},
	}
}
