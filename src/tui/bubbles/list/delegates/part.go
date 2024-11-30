package delegates

import (
	charmlist "github.com/charmbracelet/bubbles/list"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list"
)

type PartDelegate struct {
	*list.Delegate
}

var _ charmlist.ItemDelegate = &PartDelegate{}

func NewPartDelegate() *PartDelegate {
	pd := &PartDelegate{}

	pd.Delegate = list.NewDelegate(pd)

	return pd
}

func (d *PartDelegate) Title(_ charmlist.Model, _ int, item list.Item) string {
	if p, ok := item.(*puzzle.Part); ok {
		return p.Title()
	} else if i, ok := item.(charmlist.DefaultItem); ok {
		return i.Title()
	}

	return "N/A"
}

func (d *PartDelegate) Description(_ charmlist.Model, _ int, item list.Item) string {
	if p, ok := item.(*puzzle.Part); ok {
		return p.Description()
	} else if i, ok := item.(charmlist.DefaultItem); ok {
		return i.Description()
	}

	return "-"
}
