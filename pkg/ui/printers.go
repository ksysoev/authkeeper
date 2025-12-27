package ui

import (
	"fmt"
	"strings"
)

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
