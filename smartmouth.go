package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: restart (data is inconsequential, probably empty string)

func handleClientMessage(c *SpecializedClient, d []byte) {
	log.Println(string(d))
	c.hub.messages <- newMessage(c, byte(d[0]), d[1:])
}

func handleHubMessage(h *SpecializedHub, m *Message) {
	switch m.messageType {
	case byte('0'):
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(m.client.name+": "+string(m.data)))
		}
		word := strings.TrimSpace(strings.ToLower(string(m.data)))
		if len(word) >= 3 && word[0] == byte(h.start) && word[len(word)-1] == byte(h.end) {
			switch h.validWord(word) {
			case 0:
				worth := h.getWorth()
				bonus := len(word) - 2
				finalWorth := worth * bonus
				m.client.score += finalWorth
				err := h.gotIt(word)
				if err == 1 {
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte(m.client.name+" earned "+fmt.Sprint(worth)+"x"+fmt.Sprint(bonus)+"="+fmt.Sprint(finalWorth)+" points"))
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
		} else if string(m.data) == "pass" {
			m.client.pass = true
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
		} else if string(m.data) == "restart" {
			h.reset()
			for client := range h.clients {
				h.sendData(client, byte('3'), []byte(""))
				h.sendData(client, byte('0'), []byte(m.client.name+" restarted the game"))
				h.sendData(client, byte('1'), []byte(h.getScores()))
				h.sendData(client, byte('2'), []byte(h.getPrompt()))
				h.sendData(client, byte('0'), []byte("."))
				h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
			}
		}
	case byte('1'):
		name := string(m.data)
		if m.client.name == "" {
			m.client.name = name
		}
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(name+" joined"))
		}
		for client := range h.clients {
			h.sendData(client, byte('1'), []byte(h.getScores()))
		}
		h.sendData(m.client, byte('2'), []byte(h.getPrompt()))
		h.sendData(m.client, byte('0'), []byte("New letters: "+h.getPrompt()))
	}
}

func specializedInit() {
	rand.Seed(time.Now().Unix())
	buildDictionary()
}
