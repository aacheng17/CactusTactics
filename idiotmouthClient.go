package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// declaring a struct
type IdiotmouthClient struct {

	// declaring struct variable
	Client

	score int

	pass bool
}

func (c *IdiotmouthClient) handleClientMessage(d []byte) {
	log.Println(string(d))
	c.hub.Messages() <- newMessage(c, byte(d[0]), d[1:])
}

func newIdiotmouthClient(hub Hublike, conn *websocket.Conn) *IdiotmouthClient {
	ret := &IdiotmouthClient{
		Client: Client{hub: hub, conn: conn, send: make(chan []byte, 256)},
		score:  0,
		pass:   false,
	}
	ret.Client.child = ret
	return ret
}

// serveWs handles websocket requests from the peer.
func idiotmouthServeWs(hub Hublike, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newIdiotmouthClient(hub, conn)
	client.hub.Register() <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
