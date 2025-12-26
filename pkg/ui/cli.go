package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
	"golang.org/x/term"
)

// Colors
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorCyan    = "\033[36m"
	colorMagenta = "\033[35m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
)

// CoreService defines what UI needs from core (interface on consumer side)
type CoreService interface {
	AddClient(ctx context.Context, client core.Client) error
	GetClient(ctx context.Context, name string) (*core.Client, error)
	ListClients(ctx context.Context) ([]string, error)
	GetAllClients(ctx context.Context) ([]core.Client, error)
	DeleteClient(ctx context.Context, name string) error
	IssueToken(ctx context.Context, clientName string) (*core.Token, error)
	IsRepositoryInitialized() bool
	CheckPassword(ctx context.Context, password string) error
}

// CLI implements the command-line interface
type CLI struct {
	service CoreService
}

// NewCLI creates a new CLI
func NewCLI(service CoreService) *CLI {
	return &CLI{
		service: service,
	}
}

// AddClient handles the add client flow
func (c *CLI) AddClient(ctx context.Context) error {
	password, err := c.PromptMasterPassword(!c.service.IsRepositoryInitialized())
	if err != nil {
		return fmt.Errorf("failed to create master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	printTitle("🔐 Add OIDC Client")
	fmt.Println()

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
	printBox(
		fmt.Sprintf("Name:         %s", name),
		fmt.Sprintf("Client ID:    %s", clientID),
		fmt.Sprintf("Client Secret: %s", strings.Repeat("•", len(clientSecret))),
		fmt.Sprintf("Token URL:    %s", tokenURL),
		fmt.Sprintf("Scopes:       %s", strings.Join(scopes, ", ")),
	)
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
	if !c.service.IsRepositoryInitialized() {
		return fmt.Errorf("vault not initialized; please add a client first")
	}
	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	printTitle("🎫 Issue Access Token")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
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

	printBox(
		"Access Token:",
		"",
		token.AccessToken,
		"",
		fmt.Sprintf("Token Type: %s", token.TokenType),
		fmt.Sprintf("Expires In: %d seconds", token.ExpiresIn),
		fmt.Sprintf("Scope: %s", token.Scope),
	)

	fmt.Println()
	printMuted("💡 Tip: Copy the access token to use in your API requests")

	return nil
}

// ListClients handles the list clients flow
func (c *CLI) ListClients(ctx context.Context) error {
	if !c.service.IsRepositoryInitialized() {
		return fmt.Errorf("vault not initialized; please add a client first")
	}
	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	printTitle("📋 OIDC Clients")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
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
		printBox(
			fmt.Sprintf("Client ID:  %s", client.ClientID),
			fmt.Sprintf("Token URL:  %s", client.TokenURL),
			fmt.Sprintf("Scopes:     %s", strings.Join(client.Scopes, ", ")),
			fmt.Sprintf("Created:    %s", client.CreatedAt.Format("2006-01-02 15:04:05")),
		)
		fmt.Println()
	}

	return nil
}

// DeleteClient handles the delete client flow
func (c *CLI) DeleteClient(ctx context.Context) error {
	if !c.service.IsRepositoryInitialized() {
		return fmt.Errorf("vault not initialized; please add a client first")
	}
	password, err := c.PromptMasterPassword(false)
	if err != nil {
		return fmt.Errorf("failed to read master password: %w", err)
	}

	err = c.service.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	printTitle("🗑️  Delete OIDC Client")
	fmt.Println()

	// Check if repository is initialized
	if !c.service.IsRepositoryInitialized() {
		printWarning("Vault not found")
		printMuted("Use 'authkeeper add' to create vault and add your first client")
		return nil
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

// Helper functions for terminal I/O
func printTitle(text string) {
	fmt.Printf("%s%s%s%s\n", colorBold, colorMagenta, text, colorReset)
}

func printSuccess(text string) {
	fmt.Printf("%s✓ %s%s\n", colorGreen, text, colorReset)
}

func printError(text string) {
	fmt.Printf("%s✗ Error: %s%s\n", colorRed, text, colorReset)
}

func printInfo(text string) {
	fmt.Printf("%s%s%s\n", colorCyan, text, colorReset)
}

func printWarning(text string) {
	fmt.Printf("%s⚠ %s%s\n", colorYellow, text, colorReset)
}

func printMuted(text string) {
	fmt.Printf("%s%s%s\n", colorGray, text, colorReset)
}

func printProgress(text string) {
	fmt.Printf("%s%s...%s ", colorCyan, text, colorReset)
	fmt.Println()
}

func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func confirm(prompt string) bool {
	for {
		input, err := readLine(fmt.Sprintf("%s (y/n): ", prompt))
		if err != nil {
			return false
		}
		input = strings.ToLower(input)
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		printError("Please enter 'y' or 'n'")
	}
}

func selectFromList(title string, items []string) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("no items to select from")
	}

	printTitle(title)
	fmt.Println()

	for i, item := range items {
		fmt.Printf("%s%d%s. %s\n", colorCyan, i+1, colorReset, item)
	}

	fmt.Println()
	for {
		input, err := readLine("Enter number: ")
		if err != nil {
			return -1, err
		}

		var selection int
		_, err = fmt.Sscanf(input, "%d", &selection)
		if err != nil || selection < 1 || selection > len(items) {
			printError(fmt.Sprintf("Please enter a number between 1 and %d", len(items)))
			continue
		}

		return selection - 1, nil
	}
}

func printBox(lines ...string) {
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	fmt.Println("┌" + strings.Repeat("─", maxLen+2) + "┐")
	for _, line := range lines {
		padding := maxLen - len(line)
		fmt.Printf("│ %s%s │\n", line, strings.Repeat(" ", padding))
	}
	fmt.Println("└" + strings.Repeat("─", maxLen+2) + "┘")
}
