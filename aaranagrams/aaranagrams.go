package aaranagrams

import (
	"example.com/hello/core"
	u "example.com/hello/utility"
)

var name = "aaranagrams"
var Phase map[string]byte = make(map[string]byte)
var ToServerCode map[string]byte = make(map[string]byte)
var ToClientCode map[string]byte = make(map[string]byte)

type Aaranagrams struct {
	core.Game
}

func (g *Aaranagrams) Name() string {
	return name
}

func Init() core.Gamelike {
	u.GenerateEnums(name, Phase, ToServerCode, ToClientCode)
	buildDictionary()
	ret := &Aaranagrams{
		Game: *core.NewGame(NewAaranagramsHub, NewAaranagramsClient),
	}
	ret.Game.Child = ret
	return ret
}
