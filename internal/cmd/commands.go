package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ksysoev/authkeeper/internal/oauth"
	"github.com/ksysoev/authkeeper/internal/ui"
	"github.com/ksysoev/authkeeper/internal/vault"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new OIDC client to the vault",
		Long:  `Interactive command to add a new OIDC client with credentials to the encrypted vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddCommand()
		},
	}
}

func runAddCommand() error {
	ui.PrintTitle("🔐 Add OIDC Client")
	fmt.Println()

	vaultPath, err := vault.GetDefaultVaultPath()
	if err != nil {
		return fmt.Errorf("failed to get vault path: %w", err)
	}

	v := vault.New(vaultPath)

	// Check if vault exists and prompt accordingly
	isNewVault := !v.Exists()
	if isNewVault {
		ui.PrintInfo("First time setup - creating encrypted vault")
		fmt.Println()
	} else {
		ui.PrintInfo("Enter master password to unlock vault")
	}

	password, err := ui.PromptMasterPassword(isNewVault)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	// Get client details
	ui.PrintInfo("Enter client credentials")
	fmt.Println()

	name, err := ui.ReadLine("Client Name: ")
	if err != nil {
		return err
	}

	clientID, err := ui.ReadLine("Client ID: ")
	if err != nil {
		return err
	}

	clientSecret, err := ui.ReadPassword("Client Secret: ")
	if err != nil {
		return err
	}

	tokenURL, err := ui.ReadLine("Token URL: ")
	if err != nil {
		return err
	}

	scopesStr, err := ui.ReadLine("Scopes (optional, space-separated): ")
	if err != nil {
		return err
	}

	var scopes []string
	if scopesStr != "" {
		scopes = strings.Fields(scopesStr)
	}

	// Confirm
	fmt.Println()
	ui.PrintInfo("Review client details:")
	fmt.Println()
	ui.PrintBox(
		fmt.Sprintf("Name:         %s", name),
		fmt.Sprintf("Client ID:    %s", clientID),
		fmt.Sprintf("Client Secret: %s", strings.Repeat("•", len(clientSecret))),
		fmt.Sprintf("Token URL:    %s", tokenURL),
		fmt.Sprintf("Scopes:       %s", strings.Join(scopes, ", ")),
	)
	fmt.Println()

	if !ui.Confirm("Save this client?") {
		ui.PrintWarning("Cancelled")
		return nil
	}

	// Save
	spinner := ui.NewSpinner("Saving to encrypted vault")
	spinner.Start()

	client := vault.Client{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	err = v.AddClient(client, password)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	ui.PrintSuccess("Client added successfully!")
	return nil
}

func newTokenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Issue an access token",
		Long:  `Select an OIDC client from the vault and issue an access token using client credentials flow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTokenCommand()
		},
	}
}

func runTokenCommand() error {
	ui.PrintTitle("🎫 Issue Access Token")
	fmt.Println()

	vaultPath, err := vault.GetDefaultVaultPath()
	if err != nil {
		return fmt.Errorf("failed to get vault path: %w", err)
	}

	v := vault.New(vaultPath)

	// Check if vault exists
	if !v.Exists() {
		ui.PrintWarning("Vault not found")
		ui.PrintMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	// Get master password
	ui.PrintInfo("Enter master password to unlock vault")
	password, err := ui.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	// Load clients
	spinner := ui.NewSpinner("Loading clients from vault")
	spinner.Start()
	clients, err := v.ListClients(password)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	if len(clients) == 0 {
		ui.PrintWarning("No clients found in vault")
		ui.PrintMuted("Use 'authkeeper add' to add your first client")
		return nil
	}

	// Select client
	fmt.Println()
	idx, err := ui.SelectFromList("Select OIDC client:", clients)
	if err != nil {
		return err
	}

	// Get client details
	fmt.Println()
	spinner = ui.NewSpinner("Fetching access token")
	spinner.Start()

	client, err := v.GetClient(clients[idx], password)
	if err != nil {
		spinner.Stop()
		ui.PrintError(err.Error())
		return err
	}

	// Fetch token
	oauthClient := oauth.NewClient()
	token, err := oauthClient.GetToken(context.Background(), client.TokenURL, client.ClientID, client.ClientSecret, client.Scopes)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	// Display token
	ui.PrintSuccess("Token issued successfully!")
	fmt.Println()
	fmt.Printf("Client: %s\n", client.Name)
	fmt.Println()

	ui.PrintBox(
		"Access Token:",
		"",
		token.AccessToken,
		"",
		fmt.Sprintf("Token Type: %s", token.TokenType),
		fmt.Sprintf("Expires In: %d seconds", token.ExpiresIn),
		fmt.Sprintf("Scope: %s", token.Scope),
	)

	fmt.Println()
	ui.PrintMuted("💡 Tip: Copy the access token to use in your API requests")

	return nil
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all OIDC clients",
		Long:  `Display all OIDC clients stored in the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand()
		},
	}
}

func runListCommand() error {
	ui.PrintTitle("📋 OIDC Clients")
	fmt.Println()

	vaultPath, err := vault.GetDefaultVaultPath()
	if err != nil {
		return fmt.Errorf("failed to get vault path: %w", err)
	}

	v := vault.New(vaultPath)

	// Check if vault exists
	if !v.Exists() {
		ui.PrintWarning("Vault not found")
		ui.PrintMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	// Get master password
	ui.PrintInfo("Enter master password to unlock vault")
	password, err := ui.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	// Load vault
	spinner := ui.NewSpinner("Loading vault")
	spinner.Start()
	data, err := v.Load(password)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	if len(data.Clients) == 0 {
		ui.PrintWarning("No clients found")
		ui.PrintMuted("Use 'authkeeper add' to add your first client")
		return nil
	}

	fmt.Println()
	ui.PrintInfo(fmt.Sprintf("Found %d client(s)", len(data.Clients)))
	fmt.Println()

	for i, client := range data.Clients {
		fmt.Printf("%s%d. %s%s\n", ui.ColorCyan(), i+1, client.Name, ui.ColorReset())
		ui.PrintBox(
			fmt.Sprintf("Client ID:  %s", client.ClientID),
			fmt.Sprintf("Token URL:  %s", client.TokenURL),
			fmt.Sprintf("Scopes:     %s", strings.Join(client.Scopes, ", ")),
			fmt.Sprintf("Created:    %s", client.CreatedAt.Format("2006-01-02 15:04:05")),
		)
		fmt.Println()
	}

	return nil
}

func newDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete an OIDC client",
		Long:  `Select and delete an OIDC client from the vault.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteCommand()
		},
	}
}

func runDeleteCommand() error {
	ui.PrintTitle("🗑️  Delete OIDC Client")
	fmt.Println()

	vaultPath, err := vault.GetDefaultVaultPath()
	if err != nil {
		return fmt.Errorf("failed to get vault path: %w", err)
	}

	v := vault.New(vaultPath)

	// Check if vault exists
	if !v.Exists() {
		ui.PrintWarning("Vault not found")
		ui.PrintMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	// Get master password
	ui.PrintInfo("Enter master password to unlock vault")
	password, err := ui.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	// Load clients
	spinner := ui.NewSpinner("Loading clients from vault")
	spinner.Start()
	clients, err := v.ListClients(password)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	if len(clients) == 0 {
		ui.PrintWarning("No clients found in vault")
		return nil
	}

	// Select client
	fmt.Println()
	idx, err := ui.SelectFromList("Select client to delete:", clients)
	if err != nil {
		return err
	}

	// Confirm deletion
	fmt.Println()
	ui.PrintWarning(fmt.Sprintf("Are you sure you want to delete '%s'?", clients[idx]))
	ui.PrintMuted("This action cannot be undone.")
	fmt.Println()

	if !ui.Confirm("Delete this client?") {
		ui.PrintInfo("Cancelled")
		return nil
	}

	// Delete
	spinner = ui.NewSpinner("Deleting client")
	spinner.Start()
	err = v.DeleteClient(clients[idx], password)
	spinner.Stop()

	if err != nil {
		ui.PrintError(err.Error())
		return err
	}

	ui.PrintSuccess("Client deleted successfully!")
	return nil
}
