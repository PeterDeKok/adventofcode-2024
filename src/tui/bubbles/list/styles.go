package list

import (
	"github.com/charmbracelet/lipgloss"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type BoxStyles struct {
	Dimmed   lipgloss.Style
	Normal   lipgloss.Style
	Selected lipgloss.Style
}

type TextStyle struct {
	Title lipgloss.Style
	Desc  lipgloss.Style
}

type TextStyles struct {
	Dimmed   TextStyle
	Normal   TextStyle
	Selected TextStyle
}

var TitleStyle = styles.Base.
	Foreground(styles.SoftHighlightColor)

var NoItemsStyle = styles.Base.
	Foreground(styles.AppBackground)

var baseStyle = styles.Base.
	Border(lipgloss.RoundedBorder(), true, true, true, true).
	BorderForeground(styles.AppBackground).
	Foreground(styles.NormalTextColor)

var baseDelegateStyle = baseStyle.Padding(0, 1)

var (
	ComponentUnfocussedStyle = baseStyle.BorderForeground(styles.VeryDimmedColor)
	ComponentFocussedStyle   = baseStyle.BorderForeground(styles.NormalTextColor)
)

var (
	FocussedTextStyles = TextStyles{
		Dimmed: TextStyle{
			Title: styles.Base.Foreground(styles.DimmedColor).Bold(true),
			Desc:  styles.Base.Foreground(styles.DimmedColor),
		},
		Normal: TextStyle{
			Title: styles.Base.Foreground(styles.NormalTextColor).Bold(true),
			Desc:  styles.Base.Foreground(styles.NormalTextColor),
		},
		Selected: TextStyle{
			Title: styles.Base.Foreground(styles.HighlightColor).Bold(true),
			Desc:  styles.Base.Foreground(styles.SoftHighlightColor).Bold(false),
		},
	}
	UnfocussedTextStyles = TextStyles{
		Dimmed: TextStyle{
			Title: styles.Base.Foreground(styles.VeryDimmedColor),
			Desc:  styles.Base.Foreground(styles.VeryDimmedColor),
		},
		Normal: TextStyle{
			Title: styles.Base.Foreground(styles.DimmedColor),
			Desc:  styles.Base.Foreground(styles.DimmedColor),
		},
		Selected: TextStyle{
			Title: styles.Base.Foreground(styles.SoftHighlightColor),
			Desc:  styles.Base.Foreground(styles.SoftHighlightColor),
		},
	}
	FocussedBoxStyles = BoxStyles{
		Dimmed: baseDelegateStyle.Foreground(styles.DimmedColor),
		Normal: baseDelegateStyle,
		Selected: baseDelegateStyle.
			BorderForeground(styles.SoftHighlightColor).
			Foreground(styles.HighlightColor),
	}
	UnfocussedBoxStyles = BoxStyles{
		Dimmed: baseDelegateStyle.Foreground(styles.VeryDimmedColor),
		Normal: baseDelegateStyle.Foreground(styles.DimmedColor),
		Selected: baseDelegateStyle.
			BorderForeground(styles.VeryDimmedColor).
			Foreground(styles.SoftHighlightColor),
	}
	HighlightFilterRunes = lipgloss.NewStyle().Underline(true)
)
