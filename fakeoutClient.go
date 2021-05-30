package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// declaring a struct
type FakeoutClient struct {

	// declaring struct variable
	Client

	score int

	answer string

	choice int
}

func (c *FakeoutClient) handleClientMessage(d []byte) {
	log.Println(string(d))
	c.hub.Messages() <- newMessage(c, byte(d[0]), d[1:])
}

func newFakeoutClient(hub Hublike, conn *websocket.Conn) Clientlike {
	ret := &FakeoutClient{
		Client: Client{hub: hub, conn: conn, send: make(chan []byte, 256)},
		score:  0,
		answer: "",
		choice: -1,
	}
	ret.Client.child = ret
	return ret
}
