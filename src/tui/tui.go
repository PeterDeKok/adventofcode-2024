package tui2

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
)

type TuiModel struct {
	L  *log.Logger
	st *state.State

	layout *Layout
}

var _ tea.Model = &TuiModel{}

func Start(ctx context.Context, m *manage.Manager, r *remote.Client) *tea.Program {
	msgBus := bus.CreateMsgBus(ctx)

	st := state.Create(msgBus, m, r)

	tm := &TuiModel{
		L:  log.With("component", "tui"),
		st: st,

		layout: NewLayout(st, &LayoutConfig{
			Title: "Advent of Code | 2024",
		}),
	}

	program := tea.NewProgram(tm, tea.WithAltScreen())

	// Relay events to the program in the background
	go msgBus.Relay(program)

	return program
}

func (m *TuiModel) Init() tea.Cmd {
	return m.layout.Init()
}

// Update is the lifecycle phase responsible for changing the state of every
// individual component.
func (m *TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle important updates first
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, bus.KeyForceQuit) {
			return m, tea.Quit
		} else if key.Matches(msg, bus.KeyQuit) {
			m.layout.Update(msg) // Ignore any returned commands
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m, m.layout.Update(bus.ContainerSizeMsg{Width: msg.Width, Height: msg.Height})
	}

	return m, m.layout.Update(msg)
}

func (m *TuiModel) View() string {
	return m.layout.View()
}
