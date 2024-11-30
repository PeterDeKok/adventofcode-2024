package delegates

import (
	charmlist "github.com/charmbracelet/bubbles/list"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list"
)

type PuzzleDelegate struct {
	*list.Delegate
}

var _ charmlist.ItemDelegate = &PuzzleDelegate{}

func NewPuzzleDelegate() *PuzzleDelegate {
	pd := &PuzzleDelegate{}

	pd.Delegate = list.NewDelegate(pd)

	return pd
}

func (d *PuzzleDelegate) Title(_ charmlist.Model, _ int, item list.Item) string {
	if p, ok := item.(*puzzle.Puzzle); ok {
		return p.Title()
	} else if i, ok := item.(charmlist.DefaultItem); ok {
		return i.Title()
	}

	// Title
	// .Title() (.Day)

	return "N/A"
}

func (d *PuzzleDelegate) Description(_ charmlist.Model, _ int, item list.Item) string {
	if p, ok := item.(*puzzle.Puzzle); ok {
		return p.Description()
	} else if i, ok := item.(charmlist.DefaultItem); ok {
		return i.Description()
	}

	// Description
	// .Error
	// .Part1.Error
	// .Part2.Error
	// .Day -> .After(now) ? until()
	// .Part1.Loaded
	// .Part1.OK
	// .Part2.Loaded
	// .Part2.OK
	// .Part2.OK && .Part2.FastestSolution
	// .Part1.OK && .Part1.FastestSolution

	return "-"
}
