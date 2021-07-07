package fakeout

import (
	"example.com/hello/core"
	"github.com/gorilla/websocket"
)

// declaring a struct
type FakeoutClient struct {

	// declaring struct variable
	core.Client

	score int

	answer string

	choice int

	fakeouts int
}

func NewFakeoutClient(hub core.Hublike, conn *websocket.Conn) core.Clientlike {
	ret := &FakeoutClient{
		Client: core.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)},
		score:  0,
		answer: "",
		choice: -1,
	}
	ret.Client.Child = ret
	return ret
}
