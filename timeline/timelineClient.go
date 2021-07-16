package timeline

import (
	"example.com/hello/core"
	"github.com/gorilla/websocket"
)

// declaring a struct
type TimelineClient struct {

	// declaring struct variable
	core.Client

	initiative int

	hand []Event
}

func NewTimelineClient(hub core.Hublike, conn *websocket.Conn) core.Clientlike {
	ret := &TimelineClient{
		Client: core.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)},
	}
	ret.initiative = hub.(*TimelineHub).newPlayerInitiative()
	ret.Client.Child = ret
	return ret
}
