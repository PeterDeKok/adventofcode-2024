package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list/delegates"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
)

type Parts struct {
	*state.State
	*bubbles.List

	Puzzle *puzzle.Puzzle
}

func NewParts(st *state.State) *Parts {
	v := &Parts{
		State: st,

		List: bubbles.NewList("parts", []bubbles.ListItem{}, delegates.NewPartDelegate()),
	}

	v.List.Model.SetStatusBarItemName("puzzle selected", "puzzle selected")
	v.List.Model.Select(0)
	v.List.Model.SetFilteringEnabled(false)

	return v
}

func (v *Parts) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.PuzzleSelectedMsg:
		if v.Puzzle != msg.Puzzle {
			v.SetPuzzle(msg.Puzzle)
			return v.List.Update(bus.UpdateListContentMsg{})
		}

		return nil
	}

	return v.List.Update(msg)
}

func (v *Parts) Selected() *puzzle.Part {
	if i := v.List.Model.SelectedItem(); i == nil {
		return nil
	} else if p, ok := i.(*puzzle.Part); ok {
		return p
	} else {
		return nil
	}
}

func (v *Parts) SetPuzzle(p *puzzle.Puzzle) {
	v.Puzzle = p

	if p != nil {
		v.List.Model.SetItems([]bubbles.ListItem{p.Part1, p.Part2})
		v.List.Model.SetStatusBarItemName("puzzle part", "puzzle parts")
	} else {
		v.List.Model.SetItems([]bubbles.ListItem{})
		v.List.Model.SetStatusBarItemName("puzzle selected", "puzzle selected")
	}
}
