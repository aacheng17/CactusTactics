package fakeout

import "example.com/hello/core"

type Fakeout struct {
	core.Game
}

func Init() core.Gamelike {
	buildQuestions()
	return &Fakeout{
		Game: *core.NewGame("fakeout/fakeout.html", NewFakeoutHub, NewFakeoutClient),
	}
}
