package results

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/info"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
	"strings"
)

var statusColors = map[string]lipgloss.Color{
	puzzle.SolutionStatusCreated.String():        styles.VeryDimmedColor,
	puzzle.SolutionStatusCreated.String():        styles.VeryDimmedColor,
	puzzle.SolutionStatusBuildFailed.String():    styles.ErrorColor,
	puzzle.SolutionStatusBuild.String():          styles.DimmedColor,
	puzzle.SolutionStatusSamplesInvalid.String(): styles.ErrorColor,
	puzzle.SolutionStatusSamplesValid.String():   styles.DimmedColor,
	puzzle.SolutionStatusReview.String():         styles.NormalTextColor,
	puzzle.SolutionStatusInvalid.String():        styles.ErrorColor,
	puzzle.SolutionStatusValid.String():          styles.SoftHighlightColor,
}

var summaryBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(styles.VeryDimmedColor).
	Padding(1, 0)

func InfoResultToView(r *info.PartInfo, containerSize block.Sizer, containerStyle lipgloss.Style) string {
	// - Summary
	// - LastSolution
	// - Log

	refSize := containerSize.WithoutFrame(containerStyle)
	aW := refSize.W() - summaryBoxStyle.GetHorizontalFrameSize()

	// ------

	//	aSize, bSize := refSize.SplitVertical(0.25, 10, 30)

	return lipgloss.JoinVertical(lipgloss.Left,
		infoSummary(refSize, r.Summary, aW),
		infoLastSolution(refSize, r.LastSolution, aW),
		// styles.SetSize(lipgloss.NewStyle(), bSize).Render(),
	)
}

func infoSummary(size block.Sizer, summary *info.PartInfoSummary, w int) string {
	var statusColor lipgloss.Color
	var ok bool

	if statusColor, ok = statusColors[summary.BestStatus]; !ok {
		statusColor = styles.VeryDimmedColor
	}

	starSt := styles.VeryDimmedColor
	if summary.BestStatus == puzzle.SolutionStatusValid.String() {
		starSt = styles.StarColor
	}
	star := lipgloss.NewStyle().Foreground(starSt).Render("*")

	statusText := styles.RenderFancyBackground(
		lipgloss.NewStyle().
			Foreground(statusColor).
			Margin(0, 2).
			Render(star+" "+strings.ToUpper(summary.BestStatus)),
		block.Size{Width: w, Height: 3},
		styles.WithHeaderBlocks(statusColor)...,
	)

	return styles.SetSizeWithoutFrame(summaryBoxStyle, size). //.BorderForeground(statusColor)
		UnsetHeight().
		Render(lipgloss.JoinVertical(lipgloss.Left,
			statusText,
			"",
			summaryTable(summary, w, statusColor),
		))
}

func summaryTable(summary *info.PartInfoSummary, w int, statusColor lipgloss.Color) string {
	answer := ""

	if len(summary.CorrectAnswer) > 0 {
		answer = summary.CorrectAnswer
	} else if len(summary.FastestRuntime) > 0 {
		answer = "< not loaded >"
	}

	return genTable([][]string{
		{
			"solutions", or(summary.SolutionCount, "-"),
			"fastest", or(summary.FastestRuntime, "-"),
		},
		{
			"answer", or(answer, "-"),
			"finished at", or(summary.FinishedAt, "-"),
		},
	}, w, statusColor)
}

func infoLastSolution(size block.Sizer, last *info.PartInfoLastSolution, w int) string {
	color := styles.VeryDimmedColor

	if last.Status != puzzle.SolutionStatusCreated.String() && last.Status != "" {
		color = styles.DimmedColor
	}

	return styles.SetSizeWithoutFrame(summaryBoxStyle, size). //.BorderForeground(color)
		UnsetHeight().
		Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Foreground(color).Render("  Last solution"),
			"",
			lastSolutionTable(last, w, color),
		))
}

func lastSolutionTable(last *info.PartInfoLastSolution, w int, color lipgloss.Color) string {
	return genTable([][]string{
		{
			"status", or(last.Status, "-"),
			"error", or(last.Error, "-"),
		},
		{
			"runtime", or(last.Runtime, "-"),
			"finished at", or(last.FinishedAt, "-"),
		},
		{
			"answer", or(last.Answer, "-"),
			"expected", or(last.Expected, "-"),
		},
	}, w, color)
}

func genTable(data [][]string, w int, color lipgloss.Color) string {
	stylefn := func(row, col int) lipgloss.Style {
		st := lipgloss.NewStyle().Foreground(color)
		if col%2 == 0 {
			return st.Align(lipgloss.Right)
		}
		return st.Reverse(true).Align(lipgloss.Center)
	}

	return table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(stylefn).
		Width(w).
		Rows(data...).
		String()
}

func or(str, fallback string) string {
	if len(str) > 0 {
		return str
	}

	return fallback
}
