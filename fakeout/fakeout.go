package fakeout

import (
	"example.com/hello/core"
	u "example.com/hello/utility"
)

var name = "fakeout"
var Phase map[string]byte = make(map[string]byte)
var ToServerCode map[string]byte = make(map[string]byte)
var ToClientCode map[string]byte = make(map[string]byte)

type Fakeout struct {
	core.Game
}

func (g *Fakeout) Name() string {
	return name
}

func Init() core.Gamelike {
	u.GenerateEnums(name, Phase, ToServerCode, ToClientCode)
	buildQuestions()
	ret := &Fakeout{
		Game: *core.NewGame(NewFakeoutHub, NewFakeoutClient),
	}
	ret.Game.Child = ret
	return ret
}
