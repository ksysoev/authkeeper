package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#A78BFA")
	successColor   = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")
	mutedColor     = lipgloss.Color("#6B7280")
	borderColor    = lipgloss.Color("#4B5563")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1).
			MarginTop(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F3F4F6")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			MarginTop(1)

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F3F4F6")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				MarginTop(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(2)

	selectedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(secondaryColor).
				Padding(0, 3).
				MarginTop(1).
				MarginRight(2).
				Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			MarginTop(1)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	listItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			MarginLeft(2).
			MarginTop(0)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(primaryColor).
				Bold(true).
				MarginLeft(2).
				Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(2)
)

// Frames for spinner animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func getSpinnerFrame(tick int) string {
	return spinnerFrames[tick%len(spinnerFrames)]
}
