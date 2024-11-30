package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/ansi"
	"io"
)

const ellipsis = "…"

type DelegateState string

var (
	DelegateFiltering     = DelegateState("filtering")
	DelegateFilterApplied = DelegateState("filter-applied")
	DelegateSelected      = DelegateState("selected")
	DelegateNormal        = DelegateState("normal")
)

type Delegate struct {
	BoxStyles  BoxStyles
	TextStyles TextStyles

	retriever DelegateContentRetriever
}

type Delegater interface {
	list.ItemDelegate

	SetFocussed(focussed bool)
}

type DelegateContentRetriever interface {
	Title(m list.Model, index int, item list.Item) string
	Description(m list.Model, index int, item list.Item) string
}

var _ list.ItemDelegate = &Delegate{}
var _ Delegater = &Delegate{}

func NewDelegate(retriever DelegateContentRetriever) *Delegate {
	return &Delegate{
		BoxStyles:  UnfocussedBoxStyles,
		TextStyles: FocussedTextStyles,

		retriever: retriever,
	}
}

func (d *Delegate) SetFocussed(focussed bool) {
	if focussed {
		d.BoxStyles = FocussedBoxStyles
		d.TextStyles = FocussedTextStyles
	} else {
		d.BoxStyles = UnfocussedBoxStyles
		d.TextStyles = UnfocussedTextStyles
	}
}

func (d *Delegate) Height() int {
	return 4
}

func (d *Delegate) Spacing() int {
	return 0
}

func (d *Delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d *Delegate) GetHighlightedTitleAndDesc(m list.Model, index int, item list.Item) (string, string) {
	_, shouldHighlightMatches, _, textStyles := d.State(m, index)

	title := d.retriever.Title(m, index, item)
	desc := d.retriever.Description(m, index, item)

	if shouldHighlightMatches {
		matchedRunes := m.MatchesForItem(index)
		titleMatches := make([]int, 0, len(matchedRunes))
		descMatches := make([]int, 0, len(matchedRunes))
		descOffset := len(title) + 3

		for _, i := range matchedRunes {
			if i < len(title) {
				titleMatches = append(titleMatches, i)
			} else if i >= descOffset {
				descMatches = append(descMatches, i-descOffset)
			}
		}

		title = lipgloss.StyleRunes(
			title,
			titleMatches,
			textStyles.Title.Inherit(HighlightFilterRunes),
			textStyles.Title,
		)

		desc = lipgloss.StyleRunes(
			desc,
			descMatches,
			textStyles.Desc.Inherit(HighlightFilterRunes),
			textStyles.Desc,
		)
	}

	return title, desc
}

func (d *Delegate) RenderTitleAndDesc(m list.Model, index int, item list.Item) (string, string) {
	title, desc := d.GetHighlightedTitleAndDesc(m, index, item)

	_, _, _, textStyles := d.State(m, index)

	title = ansi.Truncate(
		textStyles.Title.Render(title),
		m.Width()-textStyles.Title.GetHorizontalFrameSize(),
		ellipsis,
	)
	desc = ansi.Truncate(
		textStyles.Desc.Render(desc),
		m.Width()-textStyles.Desc.GetHorizontalFrameSize(),
		ellipsis,
	)

	return title, desc
}

func (d *Delegate) RenderBox(m list.Model, index int, title, desc string) string {
	_, _, boxStyle, _ := d.State(m, index)

	return boxStyle.
		Width(m.Width() - boxStyle.GetHorizontalPadding()).
		Render(title + "\n" + desc)
}

func (d *Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	title, desc := d.RenderTitleAndDesc(m, index, item)
	itemView := d.RenderBox(m, index, title, desc)

	if _, err := fmt.Fprint(w, itemView); err != nil {
		log.Error("failed to render delegate")
	}
}

func (d *Delegate) State(m list.Model, index int) (dstate DelegateState, shouldAddMatches bool, boxstyle lipgloss.Style, textstyle TextStyle) {
	switch {
	case m.FilterState() == list.Filtering && m.FilterValue() == "":
		return DelegateFiltering, false, d.BoxStyles.Dimmed, d.TextStyles.Dimmed
	case m.FilterState() == list.Filtering:
		// ignores selected!
		return DelegateFiltering, true, d.BoxStyles.Normal, d.TextStyles.Normal
	case m.FilterState() == list.FilterApplied && index == m.Index():
		return DelegateFilterApplied, true, d.BoxStyles.Selected, d.TextStyles.Selected
	case m.FilterState() == list.FilterApplied:
		// ignores selected!
		return DelegateFilterApplied, true, d.BoxStyles.Normal, d.TextStyles.Normal
	case index == m.Index():
		return DelegateSelected, false, d.BoxStyles.Selected, d.TextStyles.Selected
	default:
		return DelegateNormal, false, d.BoxStyles.Normal, d.TextStyles.Normal
	}
}
