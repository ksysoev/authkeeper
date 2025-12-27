package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
)

// AddClient handles the add client flow
func (c *CLI) AddClient(ctx context.Context) error {
	printTitle("🔐 Add OIDC Client")
	fmt.Println()

	password, err := c.PromptMasterPassword(!c.service.IsRepositoryInitialized())
	if err != nil {
		return fmt.Errorf("failed to get master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	printInfo("Enter client credentials")
	fmt.Println()

	name, err := readLine("Client Name: ")
	if err != nil {
		return err
	}

	clientID, err := readLine("Client ID: ")
	if err != nil {
		return err
	}

	clientSecret, err := readPassword("Client Secret: ")
	if err != nil {
		return err
	}

	tokenURL, err := readLine("Token URL: ")
	if err != nil {
		return err
	}

	scopesStr, err := readLine("Scopes (optional, space-separated): ")
	if err != nil {
		return err
	}

	var scopes []string
	if scopesStr != "" {
		scopes = strings.Fields(scopesStr)
	}

	// Confirm
	fmt.Println()
	printInfo("Review client details:")
	fmt.Println()
	fmt.Printf("Name:          %s\n", name)
	fmt.Printf("Client ID:     %s\n", clientID)
	fmt.Printf("Client Secret: %s\n", strings.Repeat("•", len(clientSecret)))
	fmt.Printf("Token URL:     %s\n", tokenURL)
	fmt.Printf("Scopes:        %s\n", strings.Join(scopes, ", "))
	fmt.Println()

	if !confirm("Save this client?") {
		printWarning("Cancelled")
		return nil
	}

	// Save
	printProgress("Saving to encrypted vault")

	client := core.Client{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
		CreatedAt:    time.Now(),
	}

	err = c.service.AddClient(ctx, client)
	if err != nil {
		printError(err.Error())
		return err
	}

	printSuccess("Client added successfully!")
	return nil
}

// IssueToken handles the token issuance flow
func (c *CLI) IssueToken(ctx context.Context) error {
	printTitle("🎫 Issue Access Token")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	// Load clients
	printProgress("Loading clients from vault")
	clients, err := c.service.ListClients(ctx)
	if err != nil {
		printError(err.Error())
		return err
	}

	if len(clients) == 0 {
		printWarning("No clients found in vault")
		printMuted("Use 'authkeeper add' to add your first client")
		return nil
	}

	// Select client
	fmt.Println()
	idx, err := selectFromList("Select OIDC client:", clients)
	if err != nil {
		return err
	}

	// Fetch token
	fmt.Println()
	printProgress("Fetching access token")

	token, err := c.service.IssueToken(ctx, clients[idx])
	if err != nil {
		printError(err.Error())
		return err
	}

	// Display token
	printSuccess("Token issued successfully!")
	fmt.Println()
	fmt.Printf("Client: %s\n", clients[idx])
	fmt.Println()
	fmt.Println("Access Token:")
	fmt.Println(token.AccessToken)
	fmt.Println()
	fmt.Printf("Token Type: %s\n", token.TokenType)
	fmt.Printf("Expires In: %d seconds\n", token.ExpiresIn)
	fmt.Printf("Scope: %s\n", token.Scope)
	fmt.Println()
	printMuted("💡 Tip: Copy the access token to use in your API requests")

	return nil
}

// ListClients handles the list clients flow
func (c *CLI) ListClients(ctx context.Context) error {
	printTitle("📋 OIDC Clients")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	// Load vault
	printProgress("Loading vault")
	clients, err := c.service.GetAllClients(ctx)
	if err != nil {
		printError(err.Error())
		return err
	}

	if len(clients) == 0 {
		printWarning("No clients found")
		printMuted("Use 'authkeeper add' to add your first client")
		return nil
	}

	fmt.Println()
	printInfo(fmt.Sprintf("Found %d client(s)", len(clients)))
	fmt.Println()

	for i, client := range clients {
		fmt.Printf("%s%d. %s%s\n", colorCyan, i+1, client.Name, colorReset)
		fmt.Printf("   Client ID:  %s\n", client.ClientID)
		fmt.Printf("   Token URL:  %s\n", client.TokenURL)
		fmt.Printf("   Scopes:     %s\n", strings.Join(client.Scopes, ", "))
		fmt.Printf("   Created:    %s\n", client.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

// DeleteClient handles the delete client flow
func (c *CLI) DeleteClient(ctx context.Context) error {
	printTitle("🗑️  Delete OIDC Client")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
	}

	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	// Load clients
	printProgress("Loading clients from vault")
	clients, err := c.service.ListClients(ctx)
	if err != nil {
		printError(err.Error())
		return err
	}

	if len(clients) == 0 {
		printWarning("No clients found in vault")
		return nil
	}

	// Select client
	fmt.Println()
	idx, err := selectFromList("Select client to delete:", clients)
	if err != nil {
		return err
	}

	// Confirm deletion
	fmt.Println()
	printWarning(fmt.Sprintf("Are you sure you want to delete '%s'?", clients[idx]))
	printMuted("This action cannot be undone.")
	fmt.Println()

	if !confirm("Delete this client?") {
		printInfo("Cancelled")
		return nil
	}

	// Delete
	printProgress("Deleting client")
	err = c.service.DeleteClient(ctx, clients[idx])
	if err != nil {
		printError(err.Error())
		return err
	}

	printSuccess("Client deleted successfully!")
	return nil
}

// PromptMasterPassword prompts for master password with confirmation for new vault
func (c *CLI) PromptMasterPassword(isNewVault bool) (string, error) {
	if isNewVault {
		printTitle("🔐 Create New Vault")
		fmt.Println()
		printInfo("You're creating a new vault. Please choose a strong master password.")
		printWarning("This password encrypts all your credentials - don't forget it!")
		fmt.Println()

		for {
			password, err := readPassword("Enter master password: ")
			if err != nil {
				return "", err
			}

			if len(password) < 8 {
				printError("Password must be at least 8 characters long")
				fmt.Println()
				continue
			}

			confirm, err := readPassword("Confirm master password: ")
			if err != nil {
				return "", err
			}

			if password != confirm {
				printError("Passwords do not match. Please try again.")
				fmt.Println()
				continue
			}

			printSuccess("Master password set successfully!")
			fmt.Println()
			return password, nil
		}
	}

	// Existing vault - just ask for password once
	password, err := readPassword("Master Password: ")
	if err != nil {
		return "", err
	}
	return password, nil
}
