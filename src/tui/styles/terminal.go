package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var backupTerminalBackground = lipgloss.DefaultRenderer().Output().BackgroundColor()

func SetTerminalBackground() {
	lipgloss.DefaultRenderer().Output().SetBackgroundColor(termenv.RGBColor("#0f0f23"))
}

func ResetTerminalBackground() {
	lipgloss.DefaultRenderer().Output().SetBackgroundColor(backupTerminalBackground)
}
