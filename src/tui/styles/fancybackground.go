package styles

import "github.com/charmbracelet/lipgloss"

func RenderFancyBackground(content string, s Sizer, wsOptions ...lipgloss.WhitespaceOption) string {
	return lipgloss.Place(
		s.W(),
		max(s.H(), lipgloss.Height(content)),
		lipgloss.Center,
		lipgloss.Center,
		content,
		wsOptions...,
	)
}

func WithPaddedBlocks(color ...lipgloss.Color) []lipgloss.WhitespaceOption {
	return withBlocks("▮", color...)
}

func WithHeaderBlocks(color ...lipgloss.Color) []lipgloss.WhitespaceOption {
	return withBlocks("◎", color...)
	// ▚
	// ░
	// ▬
	// ●◎
}

func withBlocks(char string, color ...lipgloss.Color) []lipgloss.WhitespaceOption {
	c := VeryDimmedColor

	if len(color) > 0 {
		c = color[len(color)-1]
	}

	return []lipgloss.WhitespaceOption{
		lipgloss.WithWhitespaceChars(char),
		lipgloss.WithWhitespaceForeground(c),
	}
}
