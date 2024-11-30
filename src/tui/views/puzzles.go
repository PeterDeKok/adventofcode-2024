package views

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list/delegates"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
)

type Puzzles struct {
	*state.State
	*bubbles.List
}

func NewPuzzles(st *state.State) *Puzzles {
	v := &Puzzles{
		State: st,

		List: bubbles.NewList("puzzles", st.Mng.Puzzles(), delegates.NewPuzzleDelegate()),
	}

	v.List.Model.SetStatusBarItemName("puzzle", "puzzles")

	index, _, _ := st.Mng.NextPuzzle()
	v.List.Model.Select(index)

	return v
}

func (v *Puzzles) Selected() *puzzle.Puzzle {
	if i := v.List.Model.SelectedItem(); i == nil {
		return nil
	} else if p, ok := i.(*puzzle.Puzzle); ok {
		return p
	} else {
		return nil
	}
}
