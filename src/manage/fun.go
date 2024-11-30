package manage

import (
	"fmt"
	"html"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"regexp"
	"strconv"
	"strings"
)

var reTextOnly = regexp.MustCompile(`<[^>]+>`)

func ParseFunLines(r io.Reader) (string, error) {
	colors := make(map[string]string)
	lines := make([]string, 0, 27)

	state := "pre"

	for _, line := range input.LineReader(r) {
		state = determineSectionChange(line, state)

		switch state {
		case "main": // Marks the start of parsing
			continue
		case "color":
			if class, colorAnsi, ok := parseColor(line); ok {
				colors[class] = colorAnsi
			}
		case "calendar-first":
			lines = append(lines, parseCalendarFirst(line))
		case "calendar":
			if result, ok := parseCalendar(line, colors); ok {
				lines = append(lines, result)
			}
		case "calendar-end":
			return strings.Join(lines, "\n"), nil
		case "calendar-raw":
			lines = append(lines, parseCalendarRaw(line))
		}
	}

	return "", fmt.Errorf("failed to finish parsing fun lines html")
}

func determineSectionChange(line string, previousState string) (nextState string) {
	switch {
	case line == "<main>":
		return "main"
	case previousState == "pre":
		// Short circuit until the main section is reached.
		return previousState

	case strings.HasPrefix(line, ".calendar .calendar-color-"):
		return "color"

	case strings.HasPrefix(line, "<pre class=\"calendar\">"):
		return "calendar-first"

	case strings.HasPrefix(line, "<a aria-label"):
		return "calendar"

	case line == "</pre>":
		return "calendar-end"

	case previousState == "calendar":
		// The current line is NOT a normal calendar line and
		// NOT the end of the calendar section.
		// However, it should be handled. It is most likely the last static line.
		return "calendar-raw"

	default:
		return previousState
	}
}

func parseColor(line string) (class, coloransi string, found bool) {
	var colorhex string
	var ok bool

	// .calendar .calendar-color-{hash} {... color: ?#000000; ...} ...
	if line, ok = strings.CutPrefix(line, ".calendar .calendar-color-"); !ok {
		return
	}

	// {hash} {( ...;)? color: ?#000(000)?; (...; )?} ...
	if class, line, ok = strings.Cut(line, " {"); !ok {
		return
	}

	// ( ...;)? color: ?#000(000)?; (...; )?} ...
	if line, _, ok = strings.Cut(line, "; }"); !ok {
		return
	}

	// ( ...;)? color: ?#000(000)?(; ...)?
	if _, line, ok = strings.Cut(line, " color:"); !ok {
		return
	}

	line, _, _ = strings.Cut(line, ";")
	colorhex = strings.TrimSpace(line)

	if (len(colorhex) != 7 && len(colorhex) != 4) || colorhex[0] != '#' {
		return
	}

	for _, r := range colorhex[1:] {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return
		}
	}

	return class, HexToAnsi(colorhex), true
}

func parseCalendarFirst(line string) string {
	return defaultWhite(
		strings.TrimSuffix(
			strings.TrimPrefix(line, "<pre class=\"calendar\">"),
			// Removing the filler for the day-specific eol whitespace (padding, day nr & stars)
			strings.Repeat(" ", 7),
		),
	)
}

func parseCalendar(line string, colors map[string]string) (l string, found bool) {
	var ok bool

	// ^<a aria-label="Day [0-9]+(?:, (one star|two stars))?" href="/2024/day/[0-9]+" class="calendar-day([0-9]+)(?: calendar-(?:very)?complete)?">(.*)<span class="calendar-day">[0-9 ]+</span> <span class="calendar-mark-complete">\*</span><span class="calendar-mark-verycomplete">\*</span></a>$
	if _, line, ok = strings.Cut(line, ">"); !ok {
		return
	}

	if line, _, ok = strings.Cut(line, `<span class="calendar-day">`); !ok {
		return
	}

	// Close whatever color was set.
	// A massive assumption is made that open and closing tags are symetrical.
	line = strings.ReplaceAll(line, "</span>", "\033[m"+defaultColorAnsi)

	// Replace coloring span (open) tags with a reset and new escape sequence
	for hash, coloransi := range colors {
		line = strings.ReplaceAll(line, fmt.Sprintf("<span class=\"calendar-color-%s\">", hash), "\033[m"+coloransi)
	}

	// Remove any unrecognized (html) tags
	line = html.UnescapeString(reTextOnly.ReplaceAllString(line, ""))

	return defaultColor(
		strings.TrimSuffix(line, strings.Repeat(" ", 2)),
	), true
}

func parseCalendarRaw(line string) string {
	return defaultWhite(
		strings.TrimSuffix(
			line,
			// Removing the filler for the day-specific eol whitespace (padding, day nr & stars)
			strings.Repeat(" ", 7),
		),
	)
}

func defaultWhite(line string) string {
	return fmt.Sprintf("\033[38;2;204;204;204m%s\033[m", line)
}

var defaultColorAnsi = "\033[38;2;102;102;102m"

func defaultColor(line string) string {
	return fmt.Sprintf("\033[38;2;102;102;102m%s\033[m", line)
}

func HexToAnsi(hex string) string {
	if hex[0:1] == "#" {
		hex = hex[1:]
	}

	if len(hex) == 3 {
		hex = hex[0:1] + hex[0:1] + hex[1:2] + hex[1:2] + hex[2:3] + hex[2:3]
	}

	R, _ := strconv.ParseInt(hex[0:2], 16, 0)
	G, _ := strconv.ParseInt(hex[2:4], 16, 0)
	B, _ := strconv.ParseInt(hex[4:6], 16, 0)

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", R, G, B)
}
