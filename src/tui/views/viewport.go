package views

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/ansi"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/info"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/views/results"
	"strings"
	"time"
)

const HeaderLineCount = 3

type Viewport struct {
	*state.State
	*block.Block

	Part        *puzzle.Part
	Option      *op.Option
	Result      result.Result
	MainOptions bool

	Spinner       spinner.Model
	SpinnerActive bool
	Progress      progress.Model

	Table *table.Table
	VP    viewport.Model
}

type TickMsg time.Time

func NewViewport(st *state.State) *Viewport {
	b := block.New("viewport", lipgloss.NewStyle().Margin(0, 0, 0, 2))

	log.Info("Created new viewport", "style", b.Style.GetHorizontalFrameSize(), "m", b.Style.GetMarginTop())

	return &Viewport{
		State: st,
		Block: b,

		Spinner: spinner.New(spinner.WithSpinner(spinner.Points)),
		Progress: progress.New(
			progress.WithSolidFill(string(styles.SoftHighlightColor)),
		),

		Table: table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(styles.Base.Foreground(styles.VeryDimmedColor)).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return styles.Base.Foreground(styles.DimmedColor)
				} else if col == 2 {
					return styles.Base.Foreground(lipgloss.Color("#990000"))
				}

				return styles.Base
			}),
		VP: viewport.New(10, 10),
	}
}

func (v *Viewport) Init() tea.Cmd {
	return nil
}

func (v *Viewport) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}

		return nil
	case bus.PuzzlePartSelectedMsg:
		if msg.Part == nil {
			v.Part = nil
			v.Option = nil
			v.Result = nil

			cmd := v.updateProgress()
			v.updateContent()
			return tea.Batch(cmd, v.maybeInitTickers())
		}
	case bus.OptionSelectedMsg:
		v.Part = msg.Part
		v.Option = msg.Option
		v.Result = msg.Result
		v.MainOptions = msg.MainOption

		cmd := v.updateProgress()
		v.updateContent()
		return tea.Batch(cmd, v.maybeInitTickers())
	case bus.ResultUpdated:
		v.Part = msg.Part
		v.Option = msg.Option
		v.Result = msg.Result
		v.MainOptions = msg.MainOption

		cmd := v.updateProgress()
		v.updateContent()
		return tea.Batch(cmd, v.maybeInitTickers())
	case bus.ShowMainOptions:
		if !msg.Show {
			v.MainOptions = false
			v.Part = nil
			v.Option = nil
			v.Result = nil

			cmd := v.updateProgress()
			v.updateContent()
			return tea.Batch(cmd, v.maybeInitTickers())
		}

		return nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		v.Spinner, cmd = v.Spinner.Update(msg)

		v.updateContent()
		return v.maybeContinueSpinnerTick(cmd)
	case progress.FrameMsg:
		pm, cmd := v.Progress.Update(msg)
		v.Progress = pm.(progress.Model)
		v.updateContent()
		return cmd
	}

	return nil
}

func (v *Viewport) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)

	// Subtract 4 for the percentage at the end and 2 for the 'extra' margin
	v.Progress.Width = v.Size.Width - v.Style.GetHorizontalFrameSize() - 2 - 4

	v.VP.Width = v.Size.Width - v.Style.GetHorizontalFrameSize()
	v.VP.Height = v.Size.Height - v.Style.GetVerticalFrameSize() - HeaderLineCount

	v.updateContent()
}

func (v *Viewport) updateContent() {
	defer v.L.Info("updateContent", "size", v.Size)
	if (!v.MainOptions && v.Part == nil) || v.Result == nil {
		v.Content = v.Style.Render(v.getFunFallback())
		return
	}

	if opResult, ok := v.Result.OpResult().(*info.PartInfo); ok {
		v.Content = results.InfoResultToView(opResult, v.Size, v.Style)
		return
	}

	maxiw := v.Size.Width - v.Style.GetHorizontalFrameSize() - styles.ViewportInner.GetHorizontalFrameSize()
	maxiwb := v.Size.Width - v.Style.GetHorizontalFrameSize() - styles.ViewportInnerWithBorder.GetHorizontalFrameSize()
	maxiwSt := styles.ViewportInner.MaxWidth(maxiw)
	maxiwbSt := styles.ViewportInner.MaxWidth(maxiwb)

	partTitle := "Results"

	if v.Part != nil {
		partTitle = v.Part.Title()
	} else if v.Option != nil {
		partTitle = v.Option.Title()
	}

	title := lipgloss.PlaceHorizontal(
		v.Size.Width-v.Style.GetHorizontalFrameSize(),
		lipgloss.Center,
		maxiwSt.Render(partTitle),
	)

	var status string

	if err := v.Result.Error(); err != nil {
		status = "ERR"
	} else if v.Result.OK() {
		status = "OK"
	} else if v.Result.Done() {
		status = ""
	} else {
		status = v.Spinner.View()
	}

	lines := make([]string, 0, 5)

	if err := v.Result.Error(); err != nil {
		lines = append(lines, maxiwbSt.Render("\n"+styles.Base.Foreground(lipgloss.Color("#990000")).Render(ansi.Wrap(err.Error(), maxiwb, ""))))
	}

	if rows := v.Result.Rows(); len(rows) > 0 {
		// TODO content max width (and table width if > max width??)
		v.Table.ClearRows()
		v.Table.Rows(rows...)
		v.Table.Width(0)
		tableView := v.Table.Render()
		// TODO, maybe recognize separator rows and re-render a new table to separate the borders?

		if lipgloss.Width(tableView) > maxiwb {
			v.Table.Width(maxiwb)
			tableView = v.Table.Render()
		}

		lines = append(lines, tableView)
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		status,
		lipgloss.JoinVertical(lipgloss.Left,
			lines...,
		),
	)

	box := lipgloss.PlaceHorizontal(
		v.Size.Width-v.Style.GetHorizontalFrameSize(),
		lipgloss.Center,
		styles.ViewportInnerWithBorder.Render(content),
		lipgloss.WithWhitespaceChars("▮"),
		lipgloss.WithWhitespaceForeground(styles.VeryDimmedColor),
	)

	var progressView string
	if v.Result.Total() > 0 {
		v.L.Info("showpercentage", "pct", v.Progress.Percent(), "steps", v.Result.Steps(), "total", v.Result.Total())

		progressView = lipgloss.PlaceHorizontal(
			v.Size.Width-v.Style.GetHorizontalFrameSize(),
			lipgloss.Center,
			maxiwSt.Render(v.Progress.ViewAs(v.Result.Progress())),
		)
	}

	sections := make([]string, 0, 4)
	sections = append(sections, title, progressView, "")

	v.VP.SetContent(styles.RenderFancyBackground(box, block.Size{
		Width:  v.Size.Width - v.Style.GetHorizontalFrameSize(),
		Height: v.Size.Height - v.Style.GetVerticalFrameSize() - HeaderLineCount,
	}, styles.WithPaddedBlocks()...))
	v.VP.GotoBottom() // TODO Fix this to allow scrolling??
	sections = append(sections, v.VP.View())

	v.Content = v.Style.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (v *Viewport) getFunFallback() string {
	lines := []string{
		"          .-----.          .------------------.  ",
		"       .--'     |        .-'                  '-.",
		"    .--'        |        |                      |",
		".---'           |        |                      |",
		"|               |        |                      |",
		"|               |        |          ..          |",
		"|               |        |        .'  '.        |",
		"|               |        |        |    |        |",
		"|    --.        |        |        |    |        |",
		"'---'  |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        |    |        |",
		"       |        |        |        '.  .'        |",
		".------'        '------. |          ''          |",
		"|                      | |                      |",
		"|                      | |                      |",
		"|                      | |                      |",
		"|                      | '-.                  .-'",
		"'----------------------'   '------------------'  ",
	}

	fl := v.State.Mng.FunLines()

	if len(fl) > 0 {
		lines = fl
	}

	return styles.RenderFancyBackground(strings.Join(lines, "\n"), block.Size{
		Width:  v.Size.Width - v.Style.GetHorizontalFrameSize(),
		Height: v.Size.Height - v.Style.GetVerticalFrameSize(),
	})
}

func (v *Viewport) maybeInitTickers() tea.Cmd {
	if v.Result == nil || v.Result.Done() {
		v.SpinnerActive = false
		return nil
	}

	if v.SpinnerActive {
		return nil
	}

	v.SpinnerActive = true

	return v.Spinner.Tick
}

func (v *Viewport) maybeContinueSpinnerTick(cmd tea.Cmd) tea.Cmd {
	if v.Result == nil || v.Result.Done() || !v.SpinnerActive {
		return nil
	}

	return cmd
}

func (v *Viewport) updateProgress() tea.Cmd {
	//	if v.Result != nil && v.Result.Total() > 0 {
	//		return v.Progress.SetPercent(v.Result.Progress())
	//	}

	return nil
}
