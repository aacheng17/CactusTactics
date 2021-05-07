package main

import (
	"fmt"
)

var (
	startFreq = map[string]int{"a": 724, "b": 470, "c": 844, "d": 462, "e": 370, "f": 292, "g": 291, "h": 383, "i": 372, "j": 70, "k": 97, "l": 267, "m": 535, "n": 287, "o": 333, "p": 1036, "q": 49, "r": 410, "s": 1068, "t": 550, "u": 692, "v": 146, "w": 169, "x": 16, "y": 29, "z": 40}
	endFreq   = map[string]int{"a": 532, "l": 633, "i": 84, "m": 375, "k": 113, "f": 39, "n": 847, "c": 479, "e": 1874, "u": 22, "b": 20, "h": 203, "y": 1174, "s": 1102, "t": 632, "r": 656, "d": 678, "o": 97, "p": 94, "g": 273, "w": 30, "x": 34, "z": 6, "v": 2, "j": 1, "q": 0}
)

func handleClientMessage(c *Client, d []byte) {
	c.hub.messages <- newMessage(c, byte(0), d)
}

func handleHubMessage(h *Hub, m *Message) {
	switch m.messageType {
	case 0:
		for client := range h.clients {
			h.sendData(client, 0, m.data)
		}
	case 1:
		h.sendData(m.client, byte(1), []byte(fmt.Sprint(startFreq["a"])))
	}
}
