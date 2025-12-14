// Package ui provides rendering functions for the terminal user interface.
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Modern Neon color palette - vibrant, high-energy
	textColor      = lipgloss.Color("15")  // White
	accentColor    = lipgloss.Color("51")  // Bright cyan
	selectedColor  = lipgloss.Color("15")  // White (on cyan background)
	selectedBg     = lipgloss.Color("51")  // Bright cyan background
	borderColor    = lipgloss.Color("51")  // Cyan borders
	dimColor       = lipgloss.Color("244") // Gray
	highlightColor = lipgloss.Color("51")  // Cyan highlight
	liveColor      = lipgloss.Color("196") // Bright red
	goalColor      = lipgloss.Color("46")  // Bright green
	cardColor      = lipgloss.Color("226") // Bright yellow

	// Menu styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	menuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(highlightColor).
				Bold(true).
				Padding(0, 2)

	menuTitleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	menuHelpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Padding(1, 0)
)

// RenderMainMenu renders the main menu view with navigation options.
// width and height specify the terminal dimensions.
// selected indicates which menu item is currently selected (0-indexed).
// sp is the spinner model to display when loading (for other views).
// randomSpinner is the random character spinner for main view.
// loading indicates if the spinner should be shown.
func RenderMainMenu(width, height, selected int, sp spinner.Model, randomSpinner *RandomCharSpinner, loading bool) string {
	menuItems := []string{
		"Stats",
		"Live Matches",
	}

	items := make([]string, 0, len(menuItems))
	for i, item := range menuItems {
		if i == selected {
			items = append(items, menuItemSelectedStyle.Render("→ "+item))
		} else {
			items = append(items, menuItemStyle.Render("  "+item))
		}
	}

	menuContent := strings.Join(items, "\n")

	title := menuTitleStyle.Render("GOLAZO")
	help := menuHelpStyle.Render("↑/↓: navigate  Enter: select  q: quit")

	// Spinner with fixed spacing - always reserve space to prevent movement
	// Use multiple spinner instances for a longer, more prominent animation
	spinnerStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(1).
		Padding(0, 0)

	var spinnerContent string
	if loading && randomSpinner != nil {
		// Use random character spinner for main view
		spinnerContent = spinnerStyle.Render(randomSpinner.View())
	} else {
		// Reserve space even when not loading to prevent menu movement
		spinnerContent = spinnerStyle.Render("")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Repeat("\n", 1),
		spinnerContent,
		strings.Repeat("\n", 1),
		menuContent,
		strings.Repeat("\n", 2),
		help,
	)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
