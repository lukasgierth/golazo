// Package ui provides rendering functions for the terminal user interface.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// High contrast colors for minimalist theme
	textColor     = lipgloss.Color("15")  // White
	accentColor   = lipgloss.Color("3")   // Yellow
	selectedColor = lipgloss.Color("0")   // Black
	selectedBg    = lipgloss.Color("3")   // Yellow background
	borderColor   = lipgloss.Color("8")  // Gray
	dimColor      = lipgloss.Color("8")   // Gray

	// Menu styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	menuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(selectedColor).
				Background(selectedBg).
				Bold(true).
				Padding(0, 1)

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
func RenderMainMenu(width, height, selected int) string {
	menuItems := []string{
		"Live Matches",
		"Favourites",
	}

	var items []string
	for i, item := range menuItems {
		if i == selected {
			items = append(items, menuItemSelectedStyle.Render("→ "+item))
		} else {
			items = append(items, menuItemStyle.Render("  "+item))
		}
	}

	menuContent := strings.Join(items, "\n")

	title := menuTitleStyle.Render("⚽ Golazo")
	help := menuHelpStyle.Render("↑/↓: navigate  Enter: select  q: quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		strings.Repeat("\n", 2),
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

