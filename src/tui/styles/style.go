package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var Base = lipgloss.NewStyle().
	Foreground(NormalTextColor).
	Background(AppBackground)

var ViewportInner = Base.
	Padding(0, 1).
	Margin(0, 1)

var ViewportInnerWithBorder = ViewportInner.
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(NormalTextColor)
