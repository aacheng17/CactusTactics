package main

import "github.com/gorilla/websocket"

// declaring a struct
type SpecializedClient struct {

	// declaring struct variable
	Client

	score int

	pass bool
}

func newClient(hub *SpecializedHub, conn *websocket.Conn) *SpecializedClient {
	return &SpecializedClient{
		Client: Client{hub: hub, conn: conn, send: make(chan []byte, 256)},
		score:  0,
		pass:   false,
	}
}
