package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"slices"
)

type mergedKeyMap struct {
	km []help.KeyMap
}

func MergedKeyMap(km ...help.KeyMap) help.KeyMap {
	return &mergedKeyMap{
		km: km,
	}
}

func (m *mergedKeyMap) ShortHelp() []key.Binding {
	km := make([]key.Binding, 0)

	for i := len(m.km) - 1; i >= 0; i-- {
		for _, b := range m.km[i].ShortHelp() {
			// Ignore any 'earlier' (i.e. reverse range) entries
			if slices.ContainsFunc(km, func(bb key.Binding) bool {
				return slices.Equal(bb.Keys(), b.Keys())
			}) {
				continue
			}

			km = append([]key.Binding{b}, km...)
		}
	}

	return km
}

func (m *mergedKeyMap) FullHelp() [][]key.Binding {
	//TODO implement me
	panic("implement me")
}
