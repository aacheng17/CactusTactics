package main

import (
	"log"
	"net/http"

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

func newFakeoutClient(hub Hublike, conn *websocket.Conn) *FakeoutClient {
	ret := &FakeoutClient{
		Client: Client{hub: hub, conn: conn, send: make(chan []byte, 256)},
		score:  0,
		answer: "",
		choice: -1,
	}
	ret.Client.child = ret
	return ret
}

// serveWs handles websocket requests from the peer.
func fakeoutServeWs(hub Hublike, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newFakeoutClient(hub, conn)
	client.hub.Register() <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
