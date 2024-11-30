package tui2

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/views"
	"reflect"
	"sync"
)

type Layout struct {
	*block.Block
	st *state.State

	MainStyle lipgloss.Style
	MainSize  block.Size

	Title    *views.Title
	Puzzles  *views.Puzzles
	Parts    *views.Parts
	Options  *views.Options
	Viewport *views.Viewport
	Help     *views.Help
	HelpFull *views.Help
	Question *views.Question

	Focussed bus.FocusUpdate

	ShowQuestion bool
	ShowFullHelp bool
	MainOptions  bool

	opContext context.Context
	cancelOp  func()
	sync.Mutex
}

type LayoutConfig struct {
	Title string
}

func NewLayout(st *state.State, cnf *LayoutConfig) *Layout {
	b := block.New("layout")

	l := &Layout{
		Block: b,
		st:    st,

		Title:    views.NewTitle(st, cnf.Title),
		Puzzles:  views.NewPuzzles(st),
		Parts:    views.NewParts(st),
		Options:  views.NewOptions(st),
		Viewport: views.NewViewport(st),
		Help:     views.NewHelp(st),
		HelpFull: views.NewHelp(st).WithShowAll(),
		Question: views.NewQuestion(st),

		ShowFullHelp: false,
	}

	l.Focussed = l.Puzzles
	l.updateKeyMap()
	l.Puzzles.Delegate.SetFocussed(true)

	l.MainStyle = styles.Base.Margin(0, 2)

	return l
}

func (l *Layout) Init() tea.Cmd {
	return tea.Batch(
		l.Title.Init(),
		l.Puzzles.Init(),
		l.Parts.Init(),
		l.Options.Init(),
		l.Viewport.Init(),
		l.Help.Init(),
		l.HelpFull.Init(),
		l.Question.Init(),
	)
}

// Update is the lifecycle phase responsible for changing the state of every
// individual component.
// The layout component determines which signals need to go where on a macro scale.
// All expected signal types should be handled explicitly.
func (l *Layout) Update(msg tea.Msg) tea.Cmd {
	if cmd, ok := l.updateFilteringState(msg); ok {
		// Filtering updates overlap with other updates, handle them first if applicable
		return cmd
	}

	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !l.Size.Equal(msg) {
			return l.updateSizes(msg)
		}
		return nil
	case tea.KeyMsg:
		l.L.Info("keyMsg", "key", msg.String())
		if key.Matches(msg, bus.KeyQuit) {
			l.updateToAll(msg) // Ignore any returned commands

			// Returning the quit command is technically redundant,
			// but it is safe.
			return tea.Quit
		}

		if key.Matches(msg, bus.KeyQuestionCancel) && l.tryCancelQuestion() {
			return nil
		}

		if key.Matches(msg, bus.KeyCancelOp) && l.tryCancelBuildSlot() {
			return nil
		}

		if key.Matches(msg, bus.KeyShowFullHelp) && !l.ShowFullHelp {
			l.ShowFullHelp = true
			return l.updateKeyMap()
		}

		if key.Matches(msg, bus.KeyHideFullHelp) && l.ShowFullHelp {
			l.ShowFullHelp = false
			return l.updateKeyMap()
		}

		if l.ShowQuestion {
			// When a question is active,
			// it should get ownership of (almost) all key messages
			return l.Question.Update(msg)
		}

		switch l.Focussed {
		case l.Puzzles:
			if key.Matches(msg, bus.KeySelectPuzzle) {
				return tea.Batch(
					func() tea.Msg { return bus.PuzzleSelectedMsg{Puzzle: l.Puzzles.Selected()} },
					l.setFocus(l.Parts),
				)
			} else if key.Matches(msg, bus.KeyShowMainOptions) {
				l.MainOptions = true

				return tea.Batch(
					func() tea.Msg { return bus.ShowMainOptions{Show: true} },
					l.setFocus(l.Options),
				)
			}
		case l.Parts:
			if key.Matches(msg, bus.KeySelectPart) {
				return tea.Batch(
					func() tea.Msg { return bus.PuzzlePartSelectedMsg{Part: l.Parts.Selected()} },
					l.setFocus(l.Options),
				)
			}
		case l.Options:
			if key.Matches(msg, bus.KeySelectOption) {
				if opt := l.Options.Selected(); opt != nil {
					l.L.Info("option selected", "title", opt.Title())
					return l.BuildSlot(opt)
				}

				l.L.Info("no option selected")
				return nil
			}
		}

		if key.Matches(msg, bus.KeyFocusBack) {
			return l.goBack()
		}
	}

	// Propagate directed types to their intented recipients,
	// depending on the current state.
	// The layout composes the view and directs any updates depending on
	// state, like focus, flags, etc.
	switch msg := msg.(type) {
	case list.FilterMatchesMsg, tea.KeyMsg:
		return l.updateToFocussed(msg)
	case spinner.TickMsg:
		// Some functionality should also update when not focussed
		return l.updateToAll(msg)
	case progress.FrameMsg:
		// Some functionality should also update when not focussed
		return l.updateToAll(msg)
	case bus.PuzzleSelectedMsg:
		return tea.Batch(
			l.Parts.Update(msg),
			l.Options.Update(msg),
		)
	case bus.PuzzlePartSelectedMsg:
		return tea.Batch(
			l.Options.Update(msg),
			l.Viewport.Update(msg),
		)
	case bus.OptionSelectedMsg:
		// TODO We could update part & puzzle delegates??
		return tea.Batch(
			l.Options.Update(msg),
			l.Viewport.Update(msg),
		)
	case bus.ShowMainOptions:
		// TODO We could update part & puzzle delegates??
		return tea.Batch(
			l.Options.Update(msg),
			l.Viewport.Update(msg),
		)
	case bus.ResultUpdated:
		var cmds []tea.Cmd

		if msg.Result.Done() {
			if l.ShowQuestion {
				l.afterQuestion()
			}

			l.afterBuildSlot()
		} else if q := msg.Result.Question(); q != nil && q.Open() {
			cmds = append(cmds, func() tea.Msg {
				return bus.QuestionMsg{Question: q}
			})
		}

		cmds = append([]tea.Cmd{
			// TODO We could update part & puzzle delegates??
			l.Options.Update(msg),
			l.Viewport.Update(msg),
		}, cmds...)

		return tea.Batch(cmds...)
	case bus.QuestionMsg:
		if msg.Question == nil || !msg.Question.Open() {
			l.afterQuestion()
		} else {
			return l.EnableQuestion(msg)
		}

		return nil
	default:
		l.L.Error("tui update signal not recognized", "view", "layout", "signal", msg, "signaltype", reflect.TypeOf(msg))
		return l.updateToFocussed(msg)
	}
}

func (l *Layout) View() string {
	/*
	 *              TITLE
	 * +---------+---------+-------------+
	 * | PUZZLES | PARTS   | VIEWPORT    |
	 * +---------+---------+             |
	 * | OPTIONS           |             |
	 * -----------------------------------
	 *   HELP
	 */

	titleView := l.Title.View()
	puzzlesView := l.Puzzles.View()
	partsView := l.Parts.View()
	optionsView := l.Options.View()
	viewportView := l.Viewport.View()
	helpView := l.Help.View()

	main := l.MainStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				puzzlesView,
				partsView,
			),
			optionsView,
		),
		viewportView,
	))

	layout := l.Style.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleView,
		main,
		helpView,
	))

	if l.ShowQuestion {
		layout = l.Question.ViewOverlay(layout, lipgloss.Center, lipgloss.Center)
	}

	if l.ShowFullHelp {
		layout = l.HelpFull.ViewOverlay(layout, lipgloss.Bottom, lipgloss.Center)
	}

	return layout
}

func (l *Layout) updateSizes(size block.Sizeable) tea.Cmd {
	l.Size = size.Size()
	l.Style = styles.SetSizeWithoutFrame(l.Style, l.Size)

	l.MainSize = block.Size{Width: l.Width, Height: l.Height - 3 - 3}
	l.MainStyle = styles.SetSizeWithoutFrame(l.MainStyle, l.MainSize)

	/*
	 *     TITLE
	 * +---+---+---+
	 * | A | B | D |
	 * +---+---+   |
	 * |   C   |   |
	 * --------+---+
	 *     HELP
	 */

	refSize := block.Size{
		Width:  l.MainSize.Width - l.MainStyle.GetHorizontalFrameSize(),
		Height: l.MainSize.Height - l.MainStyle.GetVerticalFrameSize(),
	}

	abcSize, dSize := refSize.SplitHorizontal(0.4, 40, 40)
	abSize, cSize := abcSize.SplitVertical(0.5, 20, 20)
	aSize, bSize := abSize.SplitHorizontal(0.5, 20, 20)

	defer l.L.Info("updateSizes", "layoutsize", l.Size, "mainsize", l.MainSize)

	return tea.Batch(
		l.Title.Update(bus.ContainerSizeMsg{Width: l.MainSize.Width, Height: 3}),
		l.Puzzles.Update(bus.ContainerSizeMsg(aSize)),
		l.Parts.Update(bus.ContainerSizeMsg(bSize)),
		l.Options.Update(bus.ContainerSizeMsg(cSize)),
		l.Viewport.Update(bus.ContainerSizeMsg(dSize)),
		l.Help.Update(bus.ContainerSizeMsg{Width: l.MainSize.Width, Height: 3}),
		l.HelpFull.Update(bus.ContainerSizeMsg{Width: l.MainSize.Width - 8, Height: -1}),
		l.Question.Update(bus.ContainerSizeMsg{Width: l.MainSize.Width - 8, Height: -1}),
	)
}

func (l *Layout) updateToFocussed(msg tea.Msg) tea.Cmd {
	if l.Focussed == nil {
		return nil
	}

	return l.Focussed.Update(msg)
}

func (l *Layout) updateToAll(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 7)

	if cmd := l.Title.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Puzzles.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Parts.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Options.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Viewport.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Help.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.HelpFull.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := l.Question.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}

	return nil
}

func (l *Layout) setFocus(next bus.FocusUpdate) tea.Cmd {
	var cmds []tea.Cmd

	previous := l.Focussed
	l.Focussed = next

	if next == previous {
		return nil
	}

	switch previous {
	case l.Puzzles:
		l.Puzzles.Delegate.SetFocussed(false)
		if cmd := l.Puzzles.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case l.Parts:
		l.Parts.Delegate.SetFocussed(false)
		if cmd := l.Parts.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case l.Options:
		l.Options.Delegate.SetFocussed(false)
		if l.MainOptions {
			l.MainOptions = false
			cmds = append(cmds,
				func() tea.Msg { return bus.ShowMainOptions{Show: false} },
			)
		}
		if cmd := l.Options.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	switch next {
	case l.Puzzles:
		l.Puzzles.Delegate.SetFocussed(true)
		if cmd := l.Puzzles.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case l.Parts:
		l.Parts.Delegate.SetFocussed(true)
		if cmd := l.Parts.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case l.Options:
		l.Options.Delegate.SetFocussed(true)
		if cmd := l.Options.Update(bus.FocusChangedMsg{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	l.updateKeyMap()

	if cmd := l.Help.Update(bus.FocusChangedMsg{
		Focus: l.Focussed,
	}); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (l *Layout) goBack() tea.Cmd {
	if l.MainOptions {
		return tea.Batch(
			func() tea.Msg { return bus.PuzzlePartSelectedMsg{Part: nil} },
			func() tea.Msg { return bus.PuzzleSelectedMsg{Puzzle: nil} },
			l.setFocus(l.Puzzles),
		)
	}

	switch l.Focussed {
	case l.Options:
		return tea.Batch(
			func() tea.Msg { return bus.PuzzlePartSelectedMsg{Part: nil} },
			l.setFocus(l.Parts),
		)
	case l.Parts:
		return tea.Batch(
			func() tea.Msg { return bus.PuzzleSelectedMsg{Puzzle: nil} },
			l.setFocus(l.Puzzles),
		)
	case l.Puzzles:
		return nil
	}

	return nil
}

func (l *Layout) updateFilteringState(msg tea.Msg) (tea.Cmd, bool) {
	if l.Focussed == nil {
		return nil, false
	}

	if f, ok := l.Focussed.(bus.FocusUpdateFiltering); !ok || !f.AcceptsFilteringUpdate() {
		return nil, false
	} else if msg, ok := msg.(tea.KeyMsg); ok && key.Matches(msg, bus.KeyClearFilter, bus.KeyCancelFilter, bus.KeyAcceptFilter) {
		return f.Update(msg), true
	}

	return nil, false
}

func (l *Layout) BuildSlot(opt *op.Option) tea.Cmd {
	l.Lock()
	defer l.Unlock()

	var part *puzzle.Part
	if !l.MainOptions {
		part = l.Parts.Selected()
	}

	if l.opContext != nil {
		return func() tea.Msg {
			return bus.OptionSelectedMsg{
				Option:     opt,
				Part:       part,
				Result:     result.New().AddError(errors.New("build queue is busy")),
				MainOption: l.MainOptions,
			}
		}
	}

	p := l.Parts.Selected()
	if p == nil && !l.MainOptions {
		return func() tea.Msg {
			return bus.OptionSelectedMsg{
				Option:     opt,
				Part:       part,
				Result:     result.New().AddError(errors.New("no part selected or main options shown")),
				MainOption: l.MainOptions,
			}
		}
	}

	l.opContext, l.cancelOp = context.WithCancel(context.Background())

	l.L.Info("about to run option", "title", opt.Title())
	r := opt.Run(l.opContext, l.st.Mng, p)

	if !r.Done() {
		l.updateKeyMap()

		r.Listen(func(rr result.Result) {
			l.L.Info("result updated", "r", rr)
			l.st.Bus.Send(bus.ResultUpdated{Option: opt, Part: p, Result: r, MainOption: l.MainOptions})
		})
	}

	l.L.Info("operation running", "result", r)

	return tea.Batch(
		func() tea.Msg {
			return bus.OptionSelectedMsg{Option: opt, Part: p, Result: r, MainOption: l.MainOptions}
		},
		func() tea.Msg {
			return bus.ResultUpdated{Option: opt, Part: p, Result: r, MainOption: l.MainOptions}
		},
	)
}

func (l *Layout) tryCancelBuildSlot() bool {
	l.Lock()
	defer l.Unlock()

	if l.cancelOp == nil {
		return false
	}

	l.cancelOp()
	l.afterBuildSlot()

	return true
}

func (l *Layout) EnableQuestion(msg bus.QuestionMsg) tea.Cmd {
	l.ShowQuestion = true
	l.updateKeyMap()

	return l.Question.Update(msg)
}

func (l *Layout) tryCancelQuestion() bool {
	l.L.Info("try cancel q")
	l.Lock()
	defer l.Unlock()

	if !l.Question.Cancel() {
		l.L.Info("cancel q failed")
		return false
	}

	l.L.Info("cancel q")

	l.afterQuestion()

	return true
}

func (l *Layout) afterQuestion() {
	l.ShowQuestion = false
	l.updateKeyMap()
}

func (l *Layout) afterBuildSlot() {
	l.opContext = nil
	l.cancelOp = nil
	l.updateKeyMap()
}

func (l *Layout) updateKeyMap() tea.Cmd {
	// Separated into groups of mutualy exclusive combinations
	{
		// Esc | Go Back
		l.st.AppKeyMap.QuestionCancel.SetEnabled(l.ShowQuestion)
		l.st.AppKeyMap.CancelOp.SetEnabled(!l.ShowQuestion && l.opContext != nil)
		l.st.AppKeyMap.FocusBack.SetEnabled(!l.ShowQuestion && l.opContext == nil && l.Focussed != l.Puzzles)
	}

	{
		// Enter | Select/Show/Focus
		l.st.AppKeyMap.QuestionSelectAnswer.SetEnabled(l.ShowQuestion)

		l.st.AppKeyMap.SelectPuzzle.SetEnabled(!l.ShowQuestion && l.Focussed == l.Puzzles)
		l.st.AppKeyMap.ShowMainOptions.SetEnabled(!l.ShowQuestion && l.Focussed == l.Puzzles)
		l.st.AppKeyMap.SelectPart.SetEnabled(!l.ShowQuestion && l.Focussed == l.Parts)
		l.st.AppKeyMap.SelectPartOption.SetEnabled(!l.ShowQuestion && l.Focussed == l.Options)
	}

	{
		// tab, arrows | Navigate
		l.st.AppKeyMap.QuestionFocusNextAnswer.SetEnabled(l.ShowQuestion)
		l.st.AppKeyMap.QuestionFocusPreviousAnswer.SetEnabled(l.ShowQuestion)
	}

	{
		// ? | More help
		l.st.AppKeyMap.ShowFullHelp.SetEnabled(!l.ShowFullHelp)
		l.st.AppKeyMap.HideFullHelp.SetEnabled(l.ShowFullHelp)
	}

	return nil
}

