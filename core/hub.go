// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
)

type Hublike interface {
	GetRegister() chan Clientlike
	GetUnregister() chan Clientlike
	GetMessages() chan *Message
	HandleHubMessage(m *Message)
	Run()
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[Clientlike]bool

	// Register requests from the clients.
	Register chan Clientlike

	// Unregister requests from clients.
	Unregister chan Clientlike

	Messages chan *Message
}

func (h *Hub) RemoveClient(client Clientlike, debugMessage string) {
	delete(h.Clients, client)
	close(client.GetSend())
	log.Println(debugMessage)
}

func (h *Hub) SendData(client Clientlike, messageType byte, data []byte) {
	if len(client.GetSend()) <= cap(client.GetSend()) {
		toSend := append([]byte{messageType}, data...)
		client.GetSend() <- toSend
	} else {
		h.RemoveClient(client, "Detected and removed client with full send buffer.")
	}
}

func (h *Hub) GetRegister() chan Clientlike {
	return h.Register
}

func (h *Hub) GetUnregister() chan Clientlike {
	return h.Unregister
}

func (h *Hub) GetMessages() chan *Message {
	return h.Messages
}

func (h *Hub) HandleHubMessage(m *Message) {
	return
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				h.RemoveClient(client, "Removed client that disconnected.")
			}
		case message := <-h.Messages:
			log.Println("Received message\n\tType: " + fmt.Sprint(message.MessageType) + "\n\tData: " + string(message.Data))
			h.HandleHubMessage(message)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan Clientlike),
		Unregister: make(chan Clientlike),
		Messages:   make(chan *Message),
		Clients:    make(map[Clientlike]bool),
	}
}
