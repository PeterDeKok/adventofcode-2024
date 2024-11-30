package state

import (
	"github.com/charmbracelet/bubbles/key"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
)

type AppKeyMap struct {
	ShowFullHelp *key.Binding
	HideFullHelp *key.Binding

	Quit      *key.Binding
	ForceQuit *key.Binding

	SelectPuzzle     *key.Binding
	SelectPart       *key.Binding
	SelectPartOption *key.Binding

	ShowMainOptions *key.Binding

	FocusBack *key.Binding

	CancelOp *key.Binding

	QuestionCancel              *key.Binding
	QuestionFocusNextAnswer     *key.Binding
	QuestionFocusPreviousAnswer *key.Binding
	QuestionSelectAnswer        *key.Binding
}

func CreateAppKeyMap() *AppKeyMap {
	return &AppKeyMap{
		ShowFullHelp: &bus.KeyShowFullHelp,
		HideFullHelp: &bus.KeyHideFullHelp,

		Quit:      &bus.KeyQuit,
		ForceQuit: &bus.KeyForceQuit,

		SelectPuzzle:     &bus.KeySelectPuzzle,
		SelectPart:       &bus.KeySelectPart,
		SelectPartOption: &bus.KeySelectOption,

		ShowMainOptions: &bus.KeyShowMainOptions,

		FocusBack: &bus.KeyFocusBack,
		CancelOp:  &bus.KeyCancelOp,

		QuestionCancel:              &bus.KeyQuestionCancel,
		QuestionFocusNextAnswer:     &bus.KeyQuestionFocusNextAnswer,
		QuestionFocusPreviousAnswer: &bus.KeyQuestionFocusPreviousAnswer,
		QuestionSelectAnswer:        &bus.KeyQuestionSelectAnswer,
	}
}

func (km *AppKeyMap) KeyBindings() []key.Binding {
	return []key.Binding{
		*km.SelectPuzzle,
		*km.SelectPart,
		*km.SelectPartOption,
		*km.QuestionFocusNextAnswer,
		*km.QuestionFocusPreviousAnswer,
		*km.QuestionSelectAnswer,

		*km.FocusBack,
		*km.QuestionCancel,
		*km.CancelOp,

		*km.ShowMainOptions,
		*km.ShowFullHelp,

		*km.Quit,
		*km.ForceQuit,
	}
}

func (km *AppKeyMap) KeyBindingsFull() [][]key.Binding {
	return [][]key.Binding{
		{
			*km.SelectPuzzle,
			*km.SelectPart,
			*km.SelectPartOption,

			*km.QuestionFocusNextAnswer,
			*km.QuestionFocusPreviousAnswer,
			*km.QuestionSelectAnswer,

			*km.FocusBack,
			*km.QuestionCancel,
			*km.CancelOp,
		},
		{
			*km.ShowMainOptions,
			*km.HideFullHelp,
		},
		{
			*km.Quit,
			*km.ForceQuit,
		},
	}
}
