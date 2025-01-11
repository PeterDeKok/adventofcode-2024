package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list/delegates"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/state"
)

type Options struct {
	*state.State
	*bubbles.List

	Puzzle *puzzle.Puzzle
	Part   *puzzle.Part

	Running     bool
	MainOptions bool
}

func NewOptions(st *state.State) *Options {
	v := &Options{
		State: st,

		List: bubbles.NewList("options", []bubbles.ListItem{}, delegates.NewOptionDelegate()),
	}

	v.List.Model.SetStatusBarItemName("puzzle selected", "puzzle selected")
	v.List.Model.Select(0)

	return v
}

func (v *Options) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case bus.PuzzleSelectedMsg:
		if v.Puzzle != msg.Puzzle {
			v.SetPuzzle(msg.Puzzle)
			return v.List.Update(bus.UpdateListContentMsg{})
		}

		return nil
	case bus.PuzzlePartSelectedMsg:
		if v.Part != msg.Part {
			v.SetPart(msg.Part)
			return v.List.Update(bus.UpdateListContentMsg{})
		}

		return nil
	case bus.OptionSelectedMsg:
		v.Running = !msg.Result.Done()
		v.setItems()
		return nil
	case bus.ResultUpdated:
		v.Running = !msg.Result.Done()
		v.setItems()
	case bus.ShowMainOptions:
		if msg.Show != v.MainOptions {
			v.MainOptions = msg.Show
			v.setItems()
			return v.List.Update(bus.UpdateListContentMsg{})
		}

		return nil
	}

	return v.List.Update(msg)
}

func (v *Options) Selected() *op.Option {
	if i := v.List.Model.SelectedItem(); i == nil {
		return nil
	} else if o, ok := i.(*op.Option); ok {
		return o
	} else {
		return nil
	}
}

func (v *Options) SetPuzzle(p *puzzle.Puzzle) {
	v.Puzzle = p

	v.setItems()
}

func (v *Options) SetPart(p *puzzle.Part) {
	v.Part = p

	v.setItems()
}

func (v *Options) setItems() {
	switch {
	case v.MainOptions:
		v.setMainItems()
	case v.Puzzle == nil:
		v.setNoItems()
	case v.Part == nil:
		v.setPuzzleItems()
	default:
		v.setPartItems()
	}
}

func (v *Options) setMainItems() {
	v.List.Model.SetItems([]bubbles.ListItem{
		// op.OptionTest.SetDisabled(v.Running),
		op.OptionLoadFunLines.SetDisabled(v.Running),
		op.OptionCumulativeRuntime.SetDisabled(v.Running),
	})

	v.List.Model.SetStatusBarItemName("main option", "main options")
}

func (v *Options) setNoItems() {
	v.List.Model.SetItems([]bubbles.ListItem{})

	v.List.Model.SetStatusBarItemName("puzzle selected", "puzzle selected")
}

func (v *Options) setPuzzleItems() {
	v.List.Model.SetItems([]bubbles.ListItem{})

	v.List.Model.SetStatusBarItemName("puzzle part selected", "puzzle part selected")
}

func (v *Options) setPartItems() {
	v.List.Model.SetItems([]bubbles.ListItem{
		// op.OptionTest.SetDisabled(v.Running),
		op.OptionPartInfo.SetDisabled(v.Running),
		op.OptionValidatePart.SetDisabled(v.Running),
		op.OptionPartBoilerplate.SetDisabled(v.Running),
		op.OptionPartLoadRemote.SetDisabled(v.Running),
		op.OptionBuildPart.SetDisabled(v.Running || !v.Part.Validated),
		op.OptionRunPartSample.SetDisabled(v.Running || !v.Part.CanRunSamples()),
		op.OptionRunPart.SetDisabled(v.Running || !v.Part.CanRunInput()),
		op.OptionRecordResult.SetDisabled(v.Running || !v.Part.CanRecordResult()),
	})

	v.List.Model.SetStatusBarItemName("puzzle part option", "puzzle part options")
}
