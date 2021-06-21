package idiotmouth

import (
	"example.com/hello/core"
)

type Idiotmouth struct {
	core.Game
}

func (g *Idiotmouth) Name() string {
	return "idiotmouth"
}

func Init() core.Gamelike {
	buildDictionary()
	ret := &Idiotmouth{
		Game: *core.NewGame(NewIdiotmouthHub, NewIdiotmouthClient),
	}
	ret.Game.Child = ret
	return ret
}
