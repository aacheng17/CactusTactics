// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
)

type Hublike interface {
	Register() chan Clientlike
	Unregister() chan Clientlike
	Messages() chan *Message
	handleHubMessage(m *Message)
	run()
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[Clientlike]bool

	// Register requests from the clients.
	register chan Clientlike

	// Unregister requests from clients.
	unregister chan Clientlike

	messages chan *Message
}

func (h *Hub) removeClient(client Clientlike, debugMessage string) {
	delete(h.clients, client)
	close(client.Send())
	log.Println(debugMessage)
}

func (h *Hub) sendData(client Clientlike, messageType byte, data []byte) {
	if len(client.Send()) <= cap(client.Send()) {
		toSend := append([]byte{messageType}, data...)
		client.Send() <- toSend
	} else {
		h.removeClient(client, "Detected and removed client with full send buffer.")
	}
}

func (h *Hub) Register() chan Clientlike {
	return h.register
}

func (h *Hub) Unregister() chan Clientlike {
	return h.unregister
}

func (h *Hub) Messages() chan *Message {
	return h.messages
}

func (h *Hub) handleHubMessage(m *Message) {
	return
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.removeClient(client, "Removed client that disconnected.")
			}
		case message := <-h.messages:
			log.Println("Received message\n\tType: " + fmt.Sprint(message.messageType) + "\n\tData: " + string(message.data))
			h.handleHubMessage(message)
		}
	}
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan Clientlike),
		unregister: make(chan Clientlike),
		messages:   make(chan *Message),
		clients:    make(map[Clientlike]bool),
	}
}
