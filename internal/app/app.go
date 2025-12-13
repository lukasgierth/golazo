// Package app implements the main application model and view navigation logic.
package app

import (
	"time"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/data"
	"github.com/0xjuanma/golazo/internal/ui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	viewMain view = iota
	viewLiveMatches
)

type model struct {
	width        int
	height       int
	currentView  view
	matches      []ui.MatchDisplay
	selected     int
	matchDetails *api.MatchDetails
	liveUpdates  []string
	spinner      spinner.Model
	loading      bool
	updateGen    *data.LiveUpdateGenerator
}

// NewModel creates a new application model with default values.
func NewModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.SpinnerStyle()

	return model{
		currentView: viewMain,
		selected:    0,
		spinner:     s,
		liveUpdates: []string{},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case liveUpdateMsg:
		if m.loading {
			m.liveUpdates = append(m.liveUpdates, msg.update)
			// Continue fetching updates
			if m.updateGen != nil && m.updateGen.HasMore() {
				cmds = append(cmds, fetchLiveUpdate(m.updateGen))
			} else {
				m.loading = false
			}
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.currentView != viewMain {
				m.currentView = viewMain
				m.selected = 0
				m.matchDetails = nil
				m.liveUpdates = []string{}
				m.loading = false
				return m, nil
			}
		}

		// Handle view-specific key events
		switch m.currentView {
		case viewMain:
			return m.handleMainViewKeys(msg)
		case viewLiveMatches:
			return m.handleLiveMatchesKeys(msg)
		}
	}
	return m, nil
}

func (m model) handleMainViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < 1 {
			m.selected++
		}
		return m, nil
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "enter":
		if m.selected == 0 {
			// Stats - do nothing, stay on main view
			return m, nil
		} else if m.selected == 1 {
			// Live Matches - load matches and switch to live matches view
			matches, err := data.MockMatches()
			if err != nil {
				// If loading fails, switch view with empty matches
				// Error is silently ignored for now - could be logged in future
				m.currentView = viewLiveMatches
				m.matches = []ui.MatchDisplay{}
				return m, nil
			}

			// Convert to display format
			displayMatches := make([]ui.MatchDisplay, 0, len(matches))
			for _, match := range matches {
				displayMatches = append(displayMatches, ui.MatchDisplay{
					Match: match,
				})
			}

			m.matches = displayMatches
			m.currentView = viewLiveMatches
			m.selected = 0

			// Load details for first match if available
			if len(m.matches) > 0 {
				return m.loadMatchDetails(m.matches[0].ID)
			}

			return m, nil
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleLiveMatchesKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.matches)-1 {
			m.selected++
			// Load details for newly selected match
			if m.selected < len(m.matches) {
				return m.loadMatchDetails(m.matches[m.selected].ID)
			}
		}
		return m, nil
	case "k", "up":
		if m.selected > 0 {
			m.selected--
			// Load details for newly selected match
			if m.selected >= 0 && m.selected < len(m.matches) {
				return m.loadMatchDetails(m.matches[m.selected].ID)
			}
		}
		return m, nil
	}
	return m, nil
}

// loadMatchDetails loads match details and starts live updates.
func (m model) loadMatchDetails(matchID int) (tea.Model, tea.Cmd) {
	if details, err := data.MockMatchDetails(matchID); err == nil {
		m.matchDetails = details
		// Start live updates for selected match
		m.updateGen = data.NewLiveUpdateGenerator(matchID)
		m.liveUpdates = []string{}
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, fetchLiveUpdate(m.updateGen))
	}
	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case viewMain:
		return ui.RenderMainMenu(m.width, m.height, m.selected)
	case viewLiveMatches:
		return ui.RenderMultiPanelView(m.width, m.height, m.matches, m.selected, m.matchDetails, m.liveUpdates, m.spinner, m.loading)
	default:
		return ui.RenderMainMenu(m.width, m.height, m.selected)
	}
}

// liveUpdateMsg is a message containing a live update string.
type liveUpdateMsg struct {
	update string
}

// fetchLiveUpdate simulates fetching a live update.
func fetchLiveUpdate(gen *data.LiveUpdateGenerator) tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		if update, ok := gen.GetNextUpdate(); ok {
			return liveUpdateMsg{update: update}
		}
		return liveUpdateMsg{update: ""}
	})
}
