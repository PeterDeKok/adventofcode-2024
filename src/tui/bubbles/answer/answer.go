package answer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type Model struct {
	*block.Block

	answer *question.Answer
	selected bool
}

func New(component string, a *question.Answer) *Model {
	v := &Model{
		Block: block.New(component,
			styles.Base.
				Border(lipgloss.RoundedBorder(), false, false, true, false).
				BorderForeground(styles.VeryDimmedColor).
				Foreground(styles.NormalTextColor).
				Align(lipgloss.Center).
				Margin(0, 2),
		),

		answer: a,
		selected: false,
	}

	v.updateContent()

	return v
}

func (v *Model) Init() tea.Cmd {
	return nil
}

func (v *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}
		return nil
	case tea.KeyMsg:
		// TODO
		v.updateContent()
		return nil
	case bus.QuestionMsg:
		return v.UpdateFocus(msg.Selected == v.answer)
	default:
		log.Error("unhandled update msg", "msg", msg)
		panic("unhandled update msg")
	}
}

func (v *Model) UpdateFocus(focus bool) tea.Cmd {
	v.selected = focus

	if focus {
		v.Style = v.Style.BorderForeground(styles.HighlightColor)
	} else {
		v.Style = v.Style.BorderForeground(styles.VeryDimmedColor)
	}

	v.updateContent()

	return nil
}

func (v *Model) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)

	v.updateContent()
}

func (v *Model) updateContent() {
	mark := "  "

	if v.selected {
		mark = lipgloss.NewStyle().Foreground(styles.HighlightColor).Render("> ")
	}

	v.Content = v.Style.Render(mark + v.answer.Title)

	v.L.Info("updateContent", "size", v.Size)
}

