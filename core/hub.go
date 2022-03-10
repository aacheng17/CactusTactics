// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"

	"example.com/hello/utility"
)

type Hublike interface {
	GetRegister() chan Clientlike
	GetUnregister() chan Clientlike
	GetMessages() chan *Message
	HandleHubMessage(m *Message)
	Run()
	DisconnectClientMessage(Clientlike)
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Child Hublike
	// Registered clients.
	Clients map[Clientlike]bool

	// Register requests from the clients.
	Register chan Clientlike

	// Unregister requests from clients.
	Unregister chan Clientlike

	Messages chan *Message

	Game string

	Id string

	DeleteHubCallback func(*Hub)
}

func (h *Hub) RemoveClient(client Clientlike, debugMessage string) {
	delete(h.Clients, client)
	close(client.GetSend())
	h.Child.DisconnectClientMessage(client)
	log.Println(debugMessage)
	if len(h.Clients) == 0 {
		h.DeleteHubCallback(h)
	}
}

func (h *Hub) Broadcast(messageType byte, data []string, exceptions ...Clientlike) {
	for client := range h.Clients {
		isException := false
		for _, c := range exceptions {
			if client == c {
				isException = true
				break
			}
		}
		if !isException {
			h.SendData(client, messageType, data)
		}
	}
}

func (h *Hub) SendData(client Clientlike, messageType byte, data []string) {
	fmt.Println(fmt.Sprint("Sending message: ", rune(messageType), " ", data))
	if len(client.GetSend()) <= cap(client.GetSend()) {
		message := ""
		for _, s := range data {
			message += utility.DELIM + s
		}
		toSend := append([]byte{messageType}, []byte(message)...)
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
			log.Println(fmt.Sprint("Received message\n\tType: ", string(message.MessageCode), "\n\tData: ", message.Data))
			h.Child.HandleHubMessage(message)
		}
	}
}

func NewHub(game string, id string, deleteHubCallback func(*Hub)) *Hub {
	return &Hub{
		Register:          make(chan Clientlike),
		Unregister:        make(chan Clientlike),
		Messages:          make(chan *Message),
		Clients:           make(map[Clientlike]bool),
		Game:              game,
		Id:                id,
		DeleteHubCallback: deleteHubCallback,
	}
}
