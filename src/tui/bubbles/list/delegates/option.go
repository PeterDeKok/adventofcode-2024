package delegates

import (
	"fmt"
	charmlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type OptionDelegate struct {
	*list.Delegate
	Running *op.Option
}

var _ charmlist.ItemDelegate = &OptionDelegate{}

func NewOptionDelegate() *OptionDelegate {
	pd := &OptionDelegate{}

	pd.Delegate = list.NewDelegate(pd)

	return pd
}

func (d *OptionDelegate) Title(_ charmlist.Model, _ int, item list.Item) string {
	if i, ok := item.(*op.Option); !ok {
		return "N/A"
	} else if i.Disabled() && item == d.Running {
		return styles.Base.Foreground(styles.SoftHighlightColor).Render(fmt.Sprintf(">>> %s", i.Title()))
	} else if i.Disabled() {
		return styles.Base.Foreground(styles.DimmedColor).Render(i.Title())
	} else if item == d.Running {
		return fmt.Sprintf(">>> %s", i.Title())
	} else {
		return i.Title()
	}
}

func (d *OptionDelegate) Description(_ charmlist.Model, _ int, item list.Item) string {
	if i, ok := item.(*op.Option); !ok {
		return "-"
	} else if i.Disabled() {
		return styles.Base.Foreground(styles.DimmedColor).Render(i.Description())
	} else {
		return i.Description()
	}
}

func (d *OptionDelegate) Height() int {
	return 4
}

func (d *OptionDelegate) Update(msg tea.Msg, _ *charmlist.Model) tea.Cmd {
	switch msg := msg.(type) {
	case bus.OptionSelectedMsg:
		if msg.Result.Done() {
			d.Running = nil
		} else {
			d.Running = msg.Option
		}
	case bus.ResultUpdated:
		if msg.Result.Done() {
			d.Running = nil
		} else {
			d.Running = msg.Option
		}
	}

	return nil
}

func (d *OptionDelegate) RenderTitleAndDesc(m charmlist.Model, index int, item list.Item) (string, string) {
	title, desc := d.GetHighlightedTitleAndDesc(m, index, item)

	_, _, boxStyle, textStyles := d.State(m, index)

	w := m.Width() - textStyles.Desc.GetHorizontalFrameSize() - boxStyle.GetHorizontalFrameSize()
	title = textStyles.Title.MaxWidth(w).MaxHeight(1).Render(title)
	desc = textStyles.Desc.MaxWidth(w).MaxHeight(1).Render(desc)

	return title, desc
}

func (d *OptionDelegate) Render(w io.Writer, m charmlist.Model, index int, item list.Item) {
	_, _, boxStyle, _ := d.State(m, index)

	title, desc := d.RenderTitleAndDesc(m, index, item)
	itemView := boxStyle.
		Height(d.Height() - boxStyle.GetVerticalFrameSize()).
		MaxHeight(d.Height()).
		Width(m.Width() - boxStyle.GetHorizontalPadding()).
		Render(title + "\n" + desc)

	if _, err := fmt.Fprint(w, itemView); err != nil {
		log.Error("failed to render option delegate")
	}
}
