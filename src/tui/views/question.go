package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

const MaxWidth int = 80
var closeFocussed *struct{} = &struct{}{}

type Question struct {
	*state.State
	*block.Block

	q *question.Question

	// focus -1 for nothing
	// focus -2 for close
	// focos >= 0 for answer
	focus int

	Answers []*bubbles.Answer
}

func NewQuestion(st *state.State) *Question {
	b := block.New("question")
	if b.Size.Width > MaxWidth {
		b.Size.Width = MaxWidth
	}
	b.Style = styles.Base.
		Border(lipgloss.RoundedBorder(), true, true, true, true).
		BorderForeground(styles.SoftHighlightColor).
		MaxWidth(40).
		Padding(1, 2).
		Margin(0, 1)

	return &Question{
		State: st,
		Block: b,

		focus: -1,
	}
}

func (v *Question) Init() tea.Cmd {
	return nil
}

func (v *Question) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}
		return nil
	case bus.QuestionMsg:
		return v.setQuestion(msg.Question)
	case tea.KeyMsg:
		if key.Matches(msg, bus.KeyQuestionFocusNextAnswer) {
			return v.nextFocus()
		}

		if key.Matches(msg, bus.KeyQuestionFocusPreviousAnswer) {
			return v.previousFocus()
		}

		if key.Matches(msg, bus.KeyQuestionSelectAnswer) {
			return v.selectAnswer()
		}

		// Idea, allow answers to specify a 'hotkey'
		// 		-> Question: How to get the hotkeys into the help view
		//		   (returning an tea.Cmd is an option, but will that always properly reset??)
	}

	return nil
}

func (v *Question) Cancel() bool {
	if v.q == nil {
		return false
	}

	v.q.Cancel()

	v.q = nil
	v.setFocus(-1)

	return true
}

func (v *Question) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	if v.Size.Width > MaxWidth {
		v.Size.Width = MaxWidth
	}
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)

	s := v.Size.WithoutFrame(v.Style)

	for _, a := range v.Answers {
		a.Update(bus.ContainerSizeMsg{Width: s.Width, Height: -1})
	}

	v.updateContent()
}

func (v *Question) updateContent() tea.Cmd {
	content := ""

	if v.q != nil {
		content += v.q.Q
		content += "\n\n"

		answersContent := make([]string, 0, len(v.Answers))

		for _, a := range v.Answers {
			answersContent = append(answersContent, a.View())
		}

		content += lipgloss.JoinHorizontal(lipgloss.Bottom, answersContent...)
	}

	v.Content = v.Style.Render(content)

	v.L.Info("updateContent", "size", v.Size)

	return nil
}

func (v *Question) setQuestion(q *question.Question) tea.Cmd {
	v.L.Info("set question", "q", q)
	if v.q != q {
		v.q = q
		v.focus = -1
		v.Answers = make([]*bubbles.Answer, 0, len(q.Answers))

		for _, a := range q.Answers {
			av := bubbles.NewAnswer("answer", a)
			av.Update(bus.ContainerSizeMsg{Width: v.Size.WithoutFrame(v.Style).Width / len(q.Answers), Height: -1})

			v.Answers = append(v.Answers, av)
		}

		return v.updateContent()
	}

	return nil
}

func (v *Question) setFocus(index int) tea.Cmd {
	if v.focus != index {
		if v.focus >= 0 && v.focus < len(v.Answers) {
			v.Answers[v.focus].UpdateFocus(false)
		}

		v.focus = index

		if v.focus >= 0 && v.focus < len(v.Answers) {
			v.Answers[v.focus].UpdateFocus(true)
		}
	}

	return v.updateContent()
}

func (v *Question) nextFocus() tea.Cmd {
	var next int

	if v.focus == -2 && len(v.Answers) == 0 {
		next = -1
	} else if v.focus == -2 {
		next = 0
	} else if v.focus >= len(v.Answers) - 1 {
		// last answer focused
		next = -2
	} else {
		next = v.focus + 1
	}

	return v.setFocus(next)
}

func (v *Question) previousFocus() tea.Cmd {
	var previous int

	if v.focus == -2 && len(v.Answers) == 0 {
		previous = -1
	} else if v.focus == -2 {
		previous = len(v.Answers) - 1
	} else if v.focus == 0 {
		// first answer focused
		previous = -2
	} else {
		previous = v.focus - 1
	}

	return v.setFocus(previous)
}

func (v *Question) selectAnswer() tea.Cmd {
	if v.focus == -1 {
		return nil
	}

	if v.focus == -2 {
		v.Cancel()

		return func() tea.Msg {
			return bus.QuestionMsg{}
		}
	}

	if v.focus >= 0 && v.focus < len(v.Answers) && v.q != nil && v.focus <= len(v.q.Answers) {
		a := v.q.Answers[v.focus]
		v.q.GiveAnswer(a)
	}

	v.setFocus(-1)

	v.updateContent()

	return func() tea.Msg {
		return bus.QuestionMsg{}
	}
}
