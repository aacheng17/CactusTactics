package idiotmouth

import (
	"example.com/hello/core"
	"github.com/gorilla/websocket"
)

// declaring a struct
type IdiotmouthClient struct {

	// declaring struct variable
	core.Client

	score int

	pass bool

	highestWord string

	highestScore int
}

func NewIdiotmouthClient(hub core.Hublike, conn *websocket.Conn) core.Clientlike {
	ret := &IdiotmouthClient{
		Client: core.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)},
		score:  0,
		pass:   false,
	}
	ret.Client.Child = ret
	return ret
}
