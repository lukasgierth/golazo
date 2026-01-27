package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const standingsDialogID = "standings"

// StandingsDialog displays the league standings table for a match.
type StandingsDialog struct {
	leagueName  string
	standings   []api.LeagueTableEntry
	homeTeamID  int
	awayTeamID  int
	scrollIndex int
}

// NewStandingsDialog creates a new standings dialog.
func NewStandingsDialog(leagueName string, standings []api.LeagueTableEntry, homeTeamID, awayTeamID int) *StandingsDialog {
	return &StandingsDialog{
		leagueName:  leagueName,
		standings:   standings,
		homeTeamID:  homeTeamID,
		awayTeamID:  awayTeamID,
		scrollIndex: 0,
	}
}

// ID returns the dialog identifier.
func (d *StandingsDialog) ID() string {
	return standingsDialogID
}

// Update handles input for the standings dialog.
func (d *StandingsDialog) Update(msg tea.Msg) (Dialog, DialogAction) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "s", "q":
			return d, DialogActionClose{}
		case "j", "down":
			if d.scrollIndex < len(d.standings)-1 {
				d.scrollIndex++
			}
		case "k", "up":
			if d.scrollIndex > 0 {
				d.scrollIndex--
			}
		}
	}
	return d, nil
}

// View renders the standings table.
func (d *StandingsDialog) View(width, height int) string {
	// Calculate dialog dimensions
	dialogWidth, dialogHeight := DialogSize(width, height, 70, 25)

	// Build the table content
	content := d.renderTable(dialogWidth - 6) // Account for padding and border

	help := "j/k: scroll | esc: close"
	return RenderDialogFrameWithHelp(d.leagueName+" Standings", content, help, dialogWidth, dialogHeight)
}

// renderTable renders the standings table.
func (d *StandingsDialog) renderTable(width int) string {
	if len(d.standings) == 0 {
		return dialogDimStyle.Render("No standings data available")
	}

	var lines []string

	// Header row
	header := d.renderHeaderRow(width)
	lines = append(lines, header)

	// Separator
	separator := dialogSeparatorStyle.Render(strings.Repeat("─", width))
	lines = append(lines, separator)

	// Data rows
	for _, entry := range d.standings {
		row := d.renderTeamRow(entry, width)
		lines = append(lines, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderHeaderRow renders the table header.
func (d *StandingsDialog) renderHeaderRow(width int) string {
	posStyle := dialogHeaderStyle.Width(3).Align(lipgloss.Right)
	teamStyle := dialogHeaderStyle.Width(width - 35).Align(lipgloss.Left)
	statStyle := dialogHeaderStyle.Width(4).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		posStyle.Render("#"),
		"  ",
		teamStyle.Render("Team"),
		statStyle.Render("P"),
		statStyle.Render("W"),
		statStyle.Render("D"),
		statStyle.Render("L"),
		statStyle.Render("GD"),
		statStyle.Render("Pts"),
	)
}

// renderTeamRow renders a single team row.
func (d *StandingsDialog) renderTeamRow(entry api.LeagueTableEntry, width int) string {
	isHighlighted := entry.Team.ID == d.homeTeamID || entry.Team.ID == d.awayTeamID

	// Choose styles based on highlight
	var posStyle, teamStyle, statStyle lipgloss.Style
	if isHighlighted {
		posStyle = dialogHighlightStyle.Width(3).Align(lipgloss.Right)
		teamStyle = dialogHighlightStyle.Width(width - 35).Align(lipgloss.Left)
		statStyle = dialogHighlightStyle.Width(4).Align(lipgloss.Right)
	} else {
		posStyle = dialogValueStyle.Width(3).Align(lipgloss.Right)
		teamStyle = dialogContentStyle.Width(width - 35).Align(lipgloss.Left)
		statStyle = dialogValueStyle.Width(4).Align(lipgloss.Right)
	}

	// Truncate team name if needed
	teamName := entry.Team.ShortName
	if teamName == "" {
		teamName = entry.Team.Name
	}
	maxTeamLen := width - 38
	if len(teamName) > maxTeamLen {
		teamName = teamName[:maxTeamLen-1] + "…"
	}

	// Format goal difference with sign
	gdStr := formatGoalDifference(entry.GoalDifference)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		posStyle.Render(fmt.Sprintf("%d", entry.Position)),
		"  ",
		teamStyle.Render(teamName),
		statStyle.Render(fmt.Sprintf("%d", entry.Played)),
		statStyle.Render(fmt.Sprintf("%d", entry.Won)),
		statStyle.Render(fmt.Sprintf("%d", entry.Drawn)),
		statStyle.Render(fmt.Sprintf("%d", entry.Lost)),
		statStyle.Render(gdStr),
		statStyle.Render(fmt.Sprintf("%d", entry.Points)),
	)
}

// formatGoalDifference formats goal difference with +/- sign.
func formatGoalDifference(gd int) string {
	if gd > 0 {
		return fmt.Sprintf("+%d", gd)
	}
	return fmt.Sprintf("%d", gd)
}
