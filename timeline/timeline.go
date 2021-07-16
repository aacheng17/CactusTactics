package timeline

import (
	"example.com/hello/core"
)

type Timeline struct {
	core.Game
}

func (g *Timeline) Name() string {
	return "idiotmouth"
}

func Init() core.Gamelike {
	buildEvents()
	ret := &Timeline{
		Game: *core.NewGame(NewTimelineHub, NewTimelineClient),
	}
	ret.Game.Child = ret
	return ret
}

type Event struct {
	year  int
	title string
	info  string
}
