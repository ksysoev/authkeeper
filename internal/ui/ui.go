package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Colors
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
)

// PrintTitle prints a colored title
func PrintTitle(text string) {
	fmt.Printf("%s%s%s%s\n", colorBold, colorMagenta, text, colorReset)
}

// PrintSuccess prints a success message
func PrintSuccess(text string) {
	fmt.Printf("%s✓ %s%s\n", colorGreen, text, colorReset)
}

// PrintError prints an error message
func PrintError(text string) {
	fmt.Printf("%s✗ Error: %s%s\n", colorRed, text, colorReset)
}

// PrintInfo prints an info message
func PrintInfo(text string) {
	fmt.Printf("%s%s%s\n", colorCyan, text, colorReset)
}

// PrintWarning prints a warning message
func PrintWarning(text string) {
	fmt.Printf("%s⚠ %s%s\n", colorYellow, text, colorReset)
}

// PrintMuted prints muted text
func PrintMuted(text string) {
	fmt.Printf("%s%s%s\n", colorGray, text, colorReset)
}

// ReadLine reads a line of input from the user
func ReadLine(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ReadPassword reads a password without echoing
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}

// Confirm asks for yes/no confirmation
func Confirm(prompt string) bool {
	for {
		input, err := ReadLine(fmt.Sprintf("%s (y/n): ", prompt))
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
		PrintError("Please enter 'y' or 'n'")
	}
}

// SelectFromList presents a list and returns the selected index
func SelectFromList(title string, items []string) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("no items to select from")
	}

	PrintTitle(title)
	fmt.Println()

	for i, item := range items {
		fmt.Printf("%s%d%s. %s\n", colorCyan, i+1, colorReset, item)
	}

	fmt.Println()
	for {
		input, err := ReadLine("Enter number: ")
		if err != nil {
			return -1, err
		}

		var selection int
		_, err = fmt.Sscanf(input, "%d", &selection)
		if err != nil || selection < 1 || selection > len(items) {
			PrintError(fmt.Sprintf("Please enter a number between 1 and %d", len(items)))
			continue
		}

		return selection - 1, nil
	}
}

// PrintSeparator prints a separator line
func PrintSeparator() {
	fmt.Println(strings.Repeat("─", 60))
}

// PrintBox prints text in a box
func PrintBox(lines ...string) {
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

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// PromptMasterPassword prompts for master password with confirmation for new vault
func PromptMasterPassword(isNewVault bool) (string, error) {
	if isNewVault {
		PrintTitle("🔐 Create New Vault")
		fmt.Println()
		PrintInfo("You're creating a new vault. Please choose a strong master password.")
		PrintWarning("This password encrypts all your credentials - don't forget it!")
		fmt.Println()

		for {
			password, err := ReadPassword("Enter master password: ")
			if err != nil {
				return "", err
			}

			if len(password) < 8 {
				PrintError("Password must be at least 8 characters long")
				fmt.Println()
				continue
			}

			confirm, err := ReadPassword("Confirm master password: ")
			if err != nil {
				return "", err
			}

			if password != confirm {
				PrintError("Passwords do not match. Please try again.")
				fmt.Println()
				continue
			}

			PrintSuccess("Master password set successfully!")
			fmt.Println()
			return password, nil
		}
	}

	// Existing vault - just ask for password once
	password, err := ReadPassword("Master Password: ")
	if err != nil {
		return "", err
	}
	return password, nil
}

// Spinner shows a simple spinner animation
type Spinner struct {
	message string
	stop    chan bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan bool),
	}
}

// Start starts the spinner (for simple version, just print message)
func (s *Spinner) Start() {
	fmt.Printf("%s%s...%s ", colorCyan, s.message, colorReset)
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	fmt.Println()
}

// ColorCyan returns cyan color code
func ColorCyan() string {
	return colorCyan
}

// ColorReset returns reset color code
func ColorReset() string {
	return colorReset
}
