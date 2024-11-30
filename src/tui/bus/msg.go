package bus

import (
	tea "github.com/charmbracelet/bubbletea"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
)

type ContainerSizeMsg block.Size

func (m ContainerSizeMsg) Size() block.Size {
	return block.Size(m)
}

type PuzzleSelectedMsg struct {
	Puzzle *puzzle.Puzzle
}

type ShowMainOptions struct {
	Show bool
}

type PuzzlePartSelectedMsg struct {
	Part *puzzle.Part
}

type OptionSelectedMsg struct {
	Option     *op.Option
	Part       *puzzle.Part
	Result     result.Result
	MainOption bool
}

type UpdateListContentMsg struct{}

type FocusUpdate interface {
	Update(msg tea.Msg) tea.Cmd
}

type FocusUpdateFiltering interface {
	FocusUpdate
	AcceptsFilteringUpdate() bool
}

type FocusChangedMsg struct {
	Focus FocusUpdate
}

type ResultUpdated struct {
	Option     *op.Option
	Part       *puzzle.Part
	Result     result.Result
	MainOption bool
}

type QuestionMsg struct {
	Question *question.Question
	Selected *question.Answer
}
