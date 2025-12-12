package app

import (
	"github.com/0xjuanma/golazo/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func NewModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return ui.RenderMainView(m.width, m.height)
}
