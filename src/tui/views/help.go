package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type Help struct {
	*state.State
	*block.Block

	Model help.Model
}

func NewHelp(st *state.State) *Help {
	b := block.New("help")
	b.Style = styles.Base.
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(styles.VeryDimmedColor).
		Margin(1, 0, 1, 0).
		Padding(0, 4, 0, 4)

	return &Help{
		State: st,
		Block: b,

		Model: help.New(),
	}
}

func (v *Help) Init() tea.Cmd {
	return nil
}

func (v *Help) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}
		return nil
	case bus.FocusChangedMsg:
		v.updateContent()
	}

	return nil
}

func (v *Help) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)
	v.Model.Width = v.Width - v.Style.GetHorizontalFrameSize()

	v.updateContent()
}

func (v *Help) updateContent() {
	v.Content = v.Style.Render(v.Model.View(v))

	v.L.Info("updateContent", "size", v.Size)
}

func (v *Help) ShortHelp() []key.Binding {
	var kb []key.Binding
	kb = append(kb, v.State.AppKeyMap.KeyBindings()...)

	return kb
}

func (v *Help) FullHelp() [][]key.Binding {
	var kb [][]key.Binding
	kb = append(kb, v.State.AppKeyMap.KeyBindingsFull()...)

	return kb
}

func (v *Help) WithShowAll() *Help {
	v.Model.ShowAll = true
	v.Style = v.Style.
		Border(lipgloss.RoundedBorder(), true, true, true, true).
		BorderForeground(styles.VeryDimmedColor).
		Padding(1, 2).
		Margin(0, 0, 0, 0)

	return v
}
