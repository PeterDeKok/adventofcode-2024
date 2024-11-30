package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type Title struct {
	*state.State
	*block.Block

	Title string
}

func NewTitle(st *state.State, title string) *Title {
	b := block.New("title")
	b.Style = styles.Base.
		Foreground(styles.HighlightColor).
		AlignHorizontal(lipgloss.Center).
		Margin(1, 0, 1, 0)

	return &Title{
		State: st,
		Block: b,

		Title: title,
	}
}

func (v *Title) Init() tea.Cmd {
	return nil
}

func (v *Title) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}

		return nil
	}

	return nil
}

func (v *Title) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)

	v.updateContent()
}

func (v *Title) updateContent() {
	v.Content = v.Style.Render(v.Title)

	v.L.Info("updateContent", "size", v.Size, "title", v.Title)
}
