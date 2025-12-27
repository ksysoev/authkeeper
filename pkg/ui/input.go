package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

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
