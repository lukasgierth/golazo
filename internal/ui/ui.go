package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("6")
	secondaryColor = lipgloss.Color("8")
	successColor   = lipgloss.Color("2")
	warningColor   = lipgloss.Color("3")
	errorColor     = lipgloss.Color("1")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Align(lipgloss.Center).
			Padding(0, 1)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Align(lipgloss.Center).
			Padding(1)
)

func RenderMainView(width, height int) string {
	// Title
	title := titleStyle.Render("⚽ Golazo")

	// Subtitle
	subtitle := subtitleStyle.Render("Live football match updates in your terminal")

	// Content area
	content := containerStyle.
		Width(width - 4).
		Height(height - 10).
		Render("Welcome to Golazo!\n\nRSS feed functionality coming soon...")

	// Help text
	help := helpStyle.Render("Press 'q' to quit")

	// Combine all elements
	view := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		strings.Repeat("\n", 1),
		content,
		strings.Repeat("\n", 1),
		help,
	)

	// Center everything
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

// Helper function to truncate text to fit width
func Truncate(text string, width int) string {
	if len(text) <= width {
		return text
	}
	return text[:width-3] + "..."
}

// Helper function to wrap text
func Wrap(text string, width int) string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
