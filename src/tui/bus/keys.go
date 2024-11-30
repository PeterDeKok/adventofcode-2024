package bus

import (
	"github.com/charmbracelet/bubbles/key"
	charmlist "github.com/charmbracelet/bubbles/list"
)

var (
	KeyQuit         = key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))
	KeyForceQuit    = key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "force quit"))
	KeyShowFullHelp = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more help"))
	KeyHideFullHelp = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "hide help"))
)

var (
	KeySelectPuzzle = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select puzzle"),
		key.WithDisabled(),
	)
	KeySelectPart = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select part"),
		key.WithDisabled(),
	)
	KeySelectOption = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "run option"),
		key.WithDisabled(),
	)
	KeyShowMainOptions = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "main options"),
		key.WithDisabled(),
	)

	KeyFocusBack = key.NewBinding(
		key.WithKeys("esc", "left"),
		key.WithHelp("esc|←", "go back"),
		key.WithDisabled(),
	)
)

var (
	KeyAcceptFilter = charmlist.DefaultKeyMap().AcceptWhileFiltering
	KeyCancelFilter = charmlist.DefaultKeyMap().CancelWhileFiltering
	KeyClearFilter  = charmlist.DefaultKeyMap().ClearFilter
)

var KeyCancelOp = key.NewBinding(
	key.WithKeys("esc"),
	key.WithHelp("esc", "cancel operation"),
	key.WithDisabled(),
)

var (
	KeyQuestionCancel = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel question"),
		key.WithDisabled(),
	)
	KeyQuestionFocusNextAnswer = key.NewBinding(
		key.WithKeys("tab", "right", "down"),
		key.WithHelp("tab|→|↓", "next answer"),
		key.WithDisabled(),
	)
	KeyQuestionFocusPreviousAnswer = key.NewBinding(
		key.WithKeys("shift+tab", "left", "up"),
		key.WithHelp("shift+tab|←|↑", "previous answer"),
		key.WithDisabled(),
	)
	KeyQuestionSelectAnswer = key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter|space", "select answer"),
		key.WithDisabled(),
	)
)
