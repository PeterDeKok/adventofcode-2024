package styles

import (
	"github.com/charmbracelet/lipgloss"
)

type Sizer interface {
	W() int
	H() int
}

func SetSize(style lipgloss.Style, size Sizer) lipgloss.Style {
	return style.
		Width(size.W()).
		MaxWidth(size.W()).
		Height(size.H()).
		MaxHeight(size.H())
}

func SetSizeWithoutFrame(style lipgloss.Style, size Sizer) lipgloss.Style {
	return style.
		Width(size.W() - style.GetHorizontalMargins() - style.GetHorizontalBorderSize()).
		MaxWidth(size.W()).
		Height(size.H() - style.GetVerticalMargins() - style.GetVerticalBorderSize()).
		MaxHeight(size.H())
}
