package ui

import (
	"fmt"
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
