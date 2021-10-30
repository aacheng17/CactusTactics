package standoff

import (
	"example.com/hello/core"
)

type Standoff struct {
	core.Game
}

func (g *Standoff) Name() string {
	return "standoff"
}

func Init() core.Gamelike {
	ret := &Standoff{
		Game: *core.NewGame(NewStandoffHub, NewStandoffClient),
	}
	ret.Game.Child = ret
	return ret
}
