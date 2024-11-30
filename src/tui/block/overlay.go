package block

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"regexp"
	"strings"
)

func Overlay(overlay, onTopOf string, vAlign, hAlign lipgloss.Position, trblOffsets ...int) string {
	t, r, b, l := spreadOffsets(trblOffsets...)

	overlayLines := strings.Split(overlay, "\n")
	ow, oh := lipgloss.Width(overlay), lipgloss.Height(overlay)

	backgroundLines := strings.Split(onTopOf, "\n")
	bw, bh := lipgloss.Width(onTopOf), lipgloss.Height(onTopOf)

	var offsetLeft, offsetTop, minh int
	if ow > bw {
		offsetLeft = 0
	} else if hAlign == lipgloss.Left {
		offsetLeft = 0 + l
	} else if hAlign == lipgloss.Center {
		offsetLeft = (bw - ow + l - r) / 2
	} else if hAlign == lipgloss.Right {
		offsetLeft = bw - ow - r
	} else {
		panic("critical error: unexpected hAlign position")
	}

	if oh > bh {
		offsetTop = 0
		minh = bh
	} else if vAlign == lipgloss.Top {
		offsetTop = 0 + t
		minh = oh
	} else if vAlign == lipgloss.Center {
		offsetTop = (bh - oh + t - b) / 2
		minh = oh
	} else if vAlign == lipgloss.Bottom {
		offsetTop = bh - oh - b
		minh = oh
	} else {
		panic("critical error: unexpected vAlign position")
	}

	for y := 0; y < minh; y++ {
		overlayLine := ansi.Truncate(overlayLines[y], bw-offsetLeft, "")
		backgroundLine := backgroundLines[y+offsetTop]

		lbg := ansi.Truncate(backgroundLine, offsetLeft, "")
		backgroundParts := strings.Split(ansi.Hardwrap(
			backgroundLine,
			offsetLeft+ansi.StringWidth(overlayLine),
			true,
		), "\n")
		var rbg string

		if len(backgroundParts) >= 2 {
			if ansiStyles := regexp.MustCompile(`\x1b[[\d;]*m`).FindAllString(backgroundParts[0], -1); len(ansiStyles) > 0 {
				rbg += ansiStyles[len(ansiStyles)-1]
			}

			rbg += strings.Join(backgroundParts[1:], "")
		}

		backgroundLines[y+offsetTop] = lbg + overlayLine + rbg
	}

	return strings.Join(backgroundLines, "\n")
}

func spreadOffsets(o ...int) (int, int, int, int) {
	switch {
	case len(o) >= 4:
		return o[0], o[1], o[2], o[1]
	case len(o) == 3:
		return o[0], o[1], o[2], o[1]
	case len(o) == 2:
		return o[0], o[1], o[0], o[1]
	case len(o) == 1:
		return o[0], o[0], o[0], o[0]
	default:
		return 0, 0, 0, 0
	}
}
