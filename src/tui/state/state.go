package state

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bus"
)

type State struct {
	Bus *bus.MsgBus
	Mng *manage.Manager
	Rmt *remote.Client

	AppKeyMap *AppKeyMap
}

func Create(msgBus *bus.MsgBus, mng *manage.Manager, r *remote.Client) *State {
	app := CreateAppKeyMap()

	return &State{
		Bus: msgBus,
		Mng: mng,
		Rmt: r,

		AppKeyMap: app,
	}
}
