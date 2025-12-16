package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RenderLiveMatchesListPanel renders the left panel using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
func RenderLiveMatchesListPanel(width, height int, listModel list.Model) string {
	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelLiveMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderStatsListPanel renders the left panel for stats view using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
// Minimal design matching live view - uses list headers instead of hardcoded titles.
// List titles are only shown when there are items. Empty lists show gray messages instead.
// For 1-day view, shows both finished and upcoming lists stacked vertically.
func RenderStatsListPanel(width, height int, finishedList list.Model, upcomingList list.Model, dateRange int) string {
	// Render date range selector
	dateSelector := renderDateRangeSelector(width-6, dateRange)

	emptyStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		Padding(2, 2).
		Align(lipgloss.Center).
		Width(width - 6)

	var finishedListView string
	finishedItems := finishedList.Items()
	if len(finishedItems) == 0 {
		// No items - show empty message, no list title
		finishedListView = emptyStyle.Render(constants.EmptyNoFinishedMatches + "\n\nTry selecting a different date range (h/l keys)")
	} else {
		// Has items - show list (which includes its title)
		finishedListView = finishedList.View()
	}

	// For 1-day view, show both lists stacked vertically
	if dateRange == 1 {
		var upcomingListView string
		upcomingItems := upcomingList.Items()
		if len(upcomingItems) == 0 {
			// No upcoming matches - show empty message, no list title
			upcomingListView = emptyStyle.Render("No upcoming matches scheduled for today")
		} else {
			// Has items - show list (which includes its title)
			upcomingListView = upcomingList.View()
		}

		// Combine both lists with date selector
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			dateSelector,
			"",
			finishedListView,
			"",
			upcomingListView,
		)
		panel := panelStyle.
			Width(width).
			Height(height).
			Render(content)
		return panel
	}

	// For 3-day view, only show finished matches
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		dateSelector,
		"",
		finishedListView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// renderDateRangeSelector renders a horizontal date range selector (1d, 3d).
func renderDateRangeSelector(width int, selected int) string {
	options := []struct {
		days  int
		label string
	}{
		{1, "1d"},
		{3, "3d"},
	}

	items := make([]string, 0, len(options))
	for _, opt := range options {
		if opt.days == selected {
			// Selected option - use highlight color
			item := matchListItemSelectedStyle.Render(opt.label)
			items = append(items, item)
		} else {
			// Unselected option - use normal color
			item := matchListItemStyle.Render(opt.label)
			items = append(items, item)
		}
	}

	// Join items with separator
	separator := "  "
	selector := strings.Join(items, separator)

	// Center the selector
	selectorStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return selectorStyle.Render(selector)
}

// RenderMultiPanelViewWithList renders the live matches view with list component.
func RenderMultiPanelViewWithList(width, height int, listModel list.Model, details *api.MatchDetails, liveUpdates []string, sp spinner.Model, loading bool, randomSpinner *RandomCharSpinner, viewLoading bool) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally using style with width and alignment
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render(spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render("Loading...")
		}
	} else {
		// Reserve space with empty lines - ensure it takes up exactly spinnerHeight lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to stats view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (matches list) - shifted down
	leftPanel := RenderLiveMatchesListPanel(leftWidth, panelHeight, listModel)

	// Render right panel (match details with live updates) - shifted down
	rightPanel := renderMatchDetailsPanel(rightWidth, panelHeight, details, liveUpdates, sp, loading)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// RenderStatsViewWithList renders the stats view with list component.
// Rebuilt to match live view structure exactly: spinner at top, left panel (matches), right panel (details).
func RenderStatsViewWithList(width, height int, finishedList list.Model, upcomingList list.Model, details *api.MatchDetails, randomSpinner *RandomCharSpinner, viewLoading bool, dateRange int) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	// Match live view exactly
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space - match live view exactly
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally using style with width and alignment
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render(spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render("Loading...")
		}
	} else {
		// Reserve space with empty lines - ensure it takes up exactly spinnerHeight lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions - match live view exactly (35% left, 65% right)
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to live view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (finished matches list) - match live view structure
	// For 1-day view, combine finished and upcoming lists vertically
	leftPanel := RenderStatsListPanel(leftWidth, panelHeight, finishedList, upcomingList, dateRange)

	// Render right panel (match details) - use dedicated stats panel renderer
	rightPanel := renderStatsMatchDetailsPanel(rightWidth, panelHeight, details)

	// Create separator - match live view exactly
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels - match live view exactly
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	// Match live view exactly - use lipgloss.Left
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// renderStatsMatchDetailsPanel renders the right panel for stats view with match details.
// This is a simplified version that always shows basic match info regardless of match status.
// Designed to be expandable for more detailed statistics in the future.
func renderStatsMatchDetailsPanel(width, height int, details *api.MatchDetails) string {
	if details == nil {
		emptyMessage := lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Width(width - 6).
			PaddingTop(2).
			Render("Select a match to view details")

		return panelStyle.
			Width(width).
			Height(height).
			Render(emptyMessage)
	}

	var content strings.Builder
	infoStyle := lipgloss.NewStyle().Foreground(dimColor)

	// Match Header: Teams and Score
	teamStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Bold(true)

	homeTeam := details.HomeTeam.ShortName
	if homeTeam == "" {
		homeTeam = details.HomeTeam.Name
	}
	awayTeam := details.AwayTeam.ShortName
	if awayTeam == "" {
		awayTeam = details.AwayTeam.Name
	}

	// Score display
	if details.HomeScore != nil && details.AwayScore != nil {
		scoreText := fmt.Sprintf("%s  %d - %d  %s", homeTeam, *details.HomeScore, *details.AwayScore, awayTeam)
		content.WriteString(teamStyle.Render(scoreText))
	} else {
		matchupText := fmt.Sprintf("%s vs %s", homeTeam, awayTeam)
		content.WriteString(teamStyle.Render(matchupText))
	}
	content.WriteString("\n")

	// Status
	var statusText string
	switch details.Status {
	case api.MatchStatusFinished:
		statusText = "Full Time"
	case api.MatchStatusLive:
		if details.LiveTime != nil {
			statusText = *details.LiveTime
		} else {
			statusText = "LIVE"
		}
	case api.MatchStatusNotStarted:
		if details.MatchTime != nil {
			statusText = details.MatchTime.Format("15:04")
		} else {
			statusText = "Not Started"
		}
	default:
		statusText = string(details.Status)
	}
	content.WriteString(infoStyle.Render(statusText))
	content.WriteString("\n")

	// League
	if details.League.Name != "" {
		content.WriteString(infoStyle.Italic(true).Render(details.League.Name))
		content.WriteString("\n")
	}

	// Venue (if available)
	if details.Venue != "" {
		content.WriteString(infoStyle.Render("📍 " + details.Venue))
		content.WriteString("\n")
	}

	// Half-time score (if available)
	if details.HalfTimeScore != nil && details.HalfTimeScore.Home != nil && details.HalfTimeScore.Away != nil {
		htText := fmt.Sprintf("HT: %d - %d", *details.HalfTimeScore.Home, *details.HalfTimeScore.Away)
		content.WriteString(infoStyle.Render(htText))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Goals section (if any)
	var goals []api.MatchEvent
	for _, event := range details.Events {
		if event.Type == "goal" {
			goals = append(goals, event)
		}
	}

	if len(goals) > 0 {
		goalsTitle := lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render("⚽ Goals")
		content.WriteString(goalsTitle)
		content.WriteString("\n")

		for _, goal := range goals {
			player := "Unknown"
			if goal.Player != nil {
				player = *goal.Player
			}
			team := goal.Team.ShortName
			if team == "" {
				team = goal.Team.Name
			}
			goalLine := fmt.Sprintf("%d' %s (%s)", goal.Minute, player, team)
			content.WriteString(infoStyle.Render(goalLine))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// Cards section (count only)
	var yellowCards, redCards int
	for _, event := range details.Events {
		if event.Type == "card" {
			if event.EventType != nil {
				if *event.EventType == "yellow" {
					yellowCards++
				} else if *event.EventType == "red" {
					redCards++
				}
			}
		}
	}

	if yellowCards > 0 || redCards > 0 {
		cardsTitle := lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render("Cards")
		content.WriteString(cardsTitle)
		content.WriteString("\n")

		cardParts := make([]string, 0)
		if yellowCards > 0 {
			cardParts = append(cardParts, fmt.Sprintf("%d🟨", yellowCards))
		}
		if redCards > 0 {
			cardParts = append(cardParts, fmt.Sprintf("%d🟥", redCards))
		}
		content.WriteString(infoStyle.Render(strings.Join(cardParts, " ")))
		content.WriteString("\n")
	}

	return panelStyle.
		Width(width).
		Height(height).
		Render(content.String())
}
