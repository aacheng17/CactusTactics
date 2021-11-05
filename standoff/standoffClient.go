package standoff

import (
	"example.com/hello/core"
	"github.com/gorilla/websocket"
)

// declaring a struct
type StandoffClient struct {

	// declaring struct variable
	core.Client

	id int

	active bool

	alive bool

	decision int

	kills []string

	roundsAlive int
}

func NewStandoffClient(hub core.Hublike, conn *websocket.Conn) core.Clientlike {
	ret := &StandoffClient{
		Client:      core.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)},
		active:      false,
		alive:       true,
		decision:    -1,
		roundsAlive: -1,
	}
	ret.Client.Child = ret
	return ret
}
