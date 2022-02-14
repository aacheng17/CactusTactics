package standoff

import (
	"example.com/hello/core"
	u "example.com/hello/utility"
)

var name = "standoff"
var Phase map[string]byte = make(map[string]byte)
var ToServerCode map[string]byte = make(map[string]byte)
var ToClientCode map[string]byte = make(map[string]byte)

type Standoff struct {
	core.Game
}

func (g *Standoff) Name() string {
	return name
}

func Init() core.Gamelike {
	u.GenerateEnums(name, Phase, ToServerCode, ToClientCode)
	ret := &Standoff{
		Game: *core.NewGame(NewStandoffHub, NewStandoffClient),
	}
	ret.Game.Child = ret
	return ret
}
