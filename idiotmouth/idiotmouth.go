package idiotmouth

import "example.com/hello/core"

type Idiotmouth struct {
	core.Game
}

func Init() core.Gamelike {
	idiotmouthBuildDictionary()
	return &Idiotmouth{
		Game: *core.NewGame("idiotmouth/idiotmouth.html", NewIdiotmouthHub, NewIdiotmouthClient),
	}
}
