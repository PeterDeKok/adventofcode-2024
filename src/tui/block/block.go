package block

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type Block struct {
	Size
	Style   lipgloss.Style
	L       *log.Logger
	Content string
}

func New(component string, style ...lipgloss.Style) *Block {
	st := styles.Base

	for _, st2 := range style {
		st = st2.Inherit(st)
	}

	return &Block{
		Style: st,
		L:     log.With("component", component),
	}
}

func (bl *Block) View() string {
	return bl.Content
}

func (bl *Block) ViewOverlay(onTopOf string, vAlign, hAlign lipgloss.Position, trblOffsets ...int) string {
	return Overlay(bl.View(), onTopOf, vAlign, hAlign, trblOffsets...)
}
