package list

import (
	charmlist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/block"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
)

type List struct {
	*block.Block

	Delegate Delegater
	Model    charmlist.Model

	Component string
	Focussed  bool
}

func New[V Item](component string, items []V, delegate Delegater) *List {
	l := &List{
		Block: block.New(component, ComponentUnfocussedStyle),

		Delegate: delegate,
		Model:    charmlist.New(concreteToInterfaceList(items), delegate, 0, 0),

		Component: component,
	}

	l.Model.Title = component
	l.Model.Styles.Title = TitleStyle
	l.Model.Styles.NoItems = NoItemsStyle
	l.Model.SetShowHelp(false)
	l.Model.DisableQuitKeybindings()

	return l
}

func (v *List) Init() tea.Cmd {
	return nil
}

func (v *List) Update(msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case bus.ContainerSizeMsg:
		if !v.Size.Equal(msg) {
			v.updateSize(msg)
		}
		return nil
	case bus.UpdateListContentMsg:
		v.updateContent()
		return nil
	case bus.FocusChangedMsg:
		if msg.Focus == v && !v.Focussed {
			v.Style = v.Style.Inherit(ComponentFocussedStyle)
			v.Focussed = true
		} else if msg.Focus != v && v.Focussed {
			v.Style = v.Style.Inherit(ComponentUnfocussedStyle)
			v.Focussed = false
		}
		v.updateContent()
		return nil
	case tea.KeyMsg:
		v.Model, cmd = v.Model.Update(msg)
		v.updateContent()
		return cmd
	case charmlist.FilterMatchesMsg, spinner.TickMsg:
		v.Model, cmd = v.Model.Update(msg)
		v.updateContent()
		return cmd
	default:
		v.Model, cmd = v.Model.Update(msg)
		v.updateContent()
		return cmd
	}
}

func (v *List) updateSize(size block.Sizeable) {
	v.Size = size.Size()
	v.Style = styles.SetSizeWithoutFrame(v.Style, v.Size)
	v.Model.SetWidth(v.Width - v.Style.GetHorizontalFrameSize())
	v.Model.SetHeight(v.Height - v.Style.GetVerticalFrameSize())

	v.updateContent()
}

func (v *List) updateContent() {
	v.Content = v.Style.Render(v.Model.View())

	v.L.Info("updateContent", "size", v.Size)
}

func (v *List) AcceptsFilteringUpdate() bool {
	return v.Model.FilteringEnabled() && (v.Model.SettingFilter() || v.Model.IsFiltered())
}

func concreteToInterfaceList[V Item](v []V) []Item {
	items := make([]Item, len(v))

	for i, d := range v {
		items[i] = d
	}

	return items
}
