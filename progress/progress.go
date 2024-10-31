package progress

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 60
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type Model struct {
	Progress progress.Model
}

// NewModel creates a new progress bar model
func NewModel() Model {
	return Model{
		Progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, nil

	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.Progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.Progress.IncrPercent(0.45)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.Progress.View() + "\n\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
