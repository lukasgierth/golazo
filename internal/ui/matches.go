package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/charmbracelet/lipgloss"
)

// MatchDisplay wraps a match with display information for rendering.
type MatchDisplay struct {
	api.Match
}

var (
	// Match list styles
	matchItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	matchItemSelectedStyle = lipgloss.NewStyle().
				Foreground(selectedColor).
				Background(selectedBg).
				Bold(true).
				Padding(0, 1)

	matchHeaderStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Padding(1, 0).
				BorderBottom(true).
				BorderForeground(borderColor)

	matchScoreStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	matchLiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // Red for live
			Bold(true)

	matchTimeStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	matchLeagueStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true)
)

// RenderLiveMatches renders the live matches view with a list of matches.
// width and height specify the terminal dimensions.
// matches is the list of matches to display.
// selected indicates which match is currently selected (0-indexed).
func RenderLiveMatches(width, height int, matches []MatchDisplay, selected int) string {
	var lines []string

	// Header
	header := matchHeaderStyle.Width(width - 2).Render("Live Matches")
	lines = append(lines, header)
	lines = append(lines, "")

	if len(matches) == 0 {
		noMatches := lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No matches available")
		lines = append(lines, noMatches)
	} else {
		// Render each match
		for i, match := range matches {
			line := renderMatchItem(match, i == selected, width-4)
			lines = append(lines, line)
		}
	}

	// Help text
	help := menuHelpStyle.Render("↑/↓: navigate  Esc: back  q: quit")
	lines = append(lines, "")
	lines = append(lines, help)

	content := strings.Join(lines, "\n")
	
	// Add padding
	paddedContent := lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Left,
		lipgloss.Top,
		paddedContent,
	)
}

func renderMatchItem(match MatchDisplay, selected bool, width int) string {
	// Match status indicator
	statusIndicator := "  "
	if match.Status == api.MatchStatusLive {
		liveTime := "LIVE"
		if match.LiveTime != nil {
			liveTime = *match.LiveTime
		}
		statusIndicator = matchLiveStyle.Render("● " + liveTime)
	} else if match.Status == api.MatchStatusFinished {
		statusIndicator = matchTimeStyle.Render("FT")
	} else if match.Status == api.MatchStatusNotStarted {
		statusIndicator = matchTimeStyle.Render("VS")
	}

	// League name
	leagueName := matchLeagueStyle.Render(match.League.Name)

	// Teams and score
	homeTeam := match.HomeTeam.ShortName
	awayTeam := match.AwayTeam.ShortName

	var scoreText string
	if match.HomeScore != nil && match.AwayScore != nil {
		scoreText = matchScoreStyle.Render(
			fmt.Sprintf("%d - %d", *match.HomeScore, *match.AwayScore),
		)
	} else {
		scoreText = matchTimeStyle.Render("vs")
	}

	// Build the match line
	matchLine := fmt.Sprintf("%s  %s  %s  %s  %s",
		statusIndicator,
		leagueName,
		homeTeam,
		scoreText,
		awayTeam,
	)

	// Truncate if too long
	if len(matchLine) > width {
		matchLine = Truncate(matchLine, width)
	}

	// Apply selection style
	if selected {
		return matchItemSelectedStyle.Width(width).Render(matchLine)
	}
	return matchItemStyle.Width(width).Render(matchLine)
}

