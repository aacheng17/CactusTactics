package idiotmouth

import (
	"example.com/hello/core"
	u "example.com/hello/utility"
)

var name = "idiotmouth"
var Phase map[string]byte = make(map[string]byte)
var ToServerCode map[string]byte = make(map[string]byte)
var ToClientCode map[string]byte = make(map[string]byte)

type Idiotmouth struct {
	core.Game
}

func (g *Idiotmouth) Name() string {
	return name
}

func Init() core.Gamelike {
	u.GenerateEnums(name, Phase, ToServerCode, ToClientCode)
	buildDictionary()
	ret := &Idiotmouth{
		Game: *core.NewGame(NewIdiotmouthHub, NewIdiotmouthClient),
	}
	ret.Game.Child = ret
	return ret
}
