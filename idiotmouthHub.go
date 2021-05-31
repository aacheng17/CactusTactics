package main

import (
	"fmt"
	"log"
	"strings"
)

// declaring a struct
type IdiotmouthHub struct {

	// declaring struct variable
	Hub

	start rune

	end rune

	usedWords []string

	letters map[string]int

	wordsLeft int
}

func (h *IdiotmouthHub) getAssertedClients() map[*IdiotmouthClient]bool {
	ret := make(map[*IdiotmouthClient]bool)
	for k, v := range h.clients {
		ret[k.(*IdiotmouthClient)] = v
	}
	return ret
}

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: restart (data is inconsequential, probably empty string)

func (h *IdiotmouthHub) handleHubMessage(m *Message) {
	c := (m.client).(*IdiotmouthClient)
	switch m.messageType {
	case byte('0'):
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(c.name+": "+string(m.data)))
		}
		word := strings.TrimSpace(strings.ToLower(string(m.data)))
		if len(word) >= 3 && word[0] == byte(h.start) && word[len(word)-1] == byte(h.end) {
			switch h.validWord(word) {
			case 0:
				worth := h.getWorth()
				bonus := len(word) - 2
				finalWorth := worth * bonus
				c.score += finalWorth
				err := h.gotIt(word)
				if err == 1 {
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte(c.name+" earned "+fmt.Sprint(worth)+"x"+fmt.Sprint(bonus)+"="+fmt.Sprint(finalWorth)+" points"))
					h.sendData(client, byte('0'), []byte("."))
					h.sendData(client, byte('2'), []byte(h.getPrompt()))
					h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
					h.sendData(client, byte('1'), []byte(h.getScores()))
				}
			case 2:
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte("This word has already been used this game."))
				}
			}
		} else if string(m.data) == "/pass" {
			c.pass = true
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte("Letters passed, new letters generated"))
					h.sendData(client, byte('0'), []byte("."))
					h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
					h.sendData(client, byte('2'), []byte(h.getPrompt()))
				}
			}
		} else if string(m.data) == "/restart" {
			h.reset()
			for client := range h.clients {
				h.sendData(client, byte('3'), []byte(""))
				h.sendData(client, byte('0'), []byte(c.name+" restarted the game"))
				h.sendData(client, byte('1'), []byte(h.getScores()))
				h.sendData(client, byte('2'), []byte(h.getPrompt()))
				h.sendData(client, byte('0'), []byte("."))
				h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
			}
		}
	case byte('1'):
		name := string(m.data)
		if c.name == "" {
			c.name = name
		}
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(name+" joined"))
		}
		for client := range h.clients {
			h.sendData(client, byte('1'), []byte(h.getScores()))
		}
		h.sendData(c, byte('2'), []byte(h.getPrompt()))
		h.sendData(c, byte('0'), []byte("New letters: "+h.getPrompt()))
	}
}

func (h *IdiotmouthHub) run() {
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

func newIdiotmouthHub() Hublike {
	h := &IdiotmouthHub{
		Hub:     *newHub(),
		letters: make(map[string]int),
	}
	h.reset()
	return h
}
