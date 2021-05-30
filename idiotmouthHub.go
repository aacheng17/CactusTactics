package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
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

func (h *IdiotmouthHub) validWord(str string) int {
	for _, v := range h.usedWords {
		if v == str {
			return 2
		}
	}
	for _, v := range words {
		if v == str {
			return 0
		}
	}
	return 1
}

func (h *IdiotmouthHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.pass = false
	}
	h.wordsLeft = len(words)
	for k, v := range letters {
		h.letters[k] = v
	}
	h.genNextLetters()
}

func (h *IdiotmouthHub) resetPass() {
	for client := range h.getAssertedClients() {
		client.pass = false
	}
}

func (h *IdiotmouthHub) pass() int {
	h.resetPass()
	h.wordsLeft -= h.letters[string(h.start)+string(h.end)]
	h.letters[string(h.start)+string(h.end)] = 0
	return h.genNextLetters()
}

func (h *IdiotmouthHub) getMajorityPass() bool {
	count := 0
	for client := range h.getAssertedClients() {
		if client.pass {
			count++
		}
	}
	return count*2 > len(h.clients)
}

func (h *IdiotmouthHub) gotIt(word string) int {
	h.resetPass()
	h.usedWords = append(h.usedWords, word)
	h.wordsLeft--
	h.letters[string(h.start)+string(h.end)]--
	return h.genNextLetters()
}

func (h *IdiotmouthHub) genNextLetters() int {
	if h.wordsLeft <= 0 {
		return 1
	}
	r := rand.Intn(h.wordsLeft)
	c := 0
	for lets, freq := range h.letters {
		c += freq
		if r < c {
			h.start = rune(lets[0])
			h.end = rune(lets[1])
			break
		}
	}
	return 0
}

func (h *IdiotmouthHub) getWorth() int {
	return int(50-50*(float32(letters[string(h.start)+string(h.end)]-minFreq)/float32(maxFreq-minFreq))) + 50
}

func (h *IdiotmouthHub) getPrompt() string {
	ret := string(h.start) + "*" + string(h.end)
	bonus := h.getWorth()
	ret += ", worth " + fmt.Sprint(bonus) + " points. There are " + fmt.Sprint(h.letters[string(h.start)+string(h.end)]) + " possible words"
	return ret
}

func (h *IdiotmouthHub) getScores() string {
	keys := make([]*IdiotmouthClient, 0, len(h.clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	scores := ""
	for _, client := range keys {
		if client.name == "" {
			continue
		}
		scores += client.name + ": " + fmt.Sprint(client.score) + "; "
	}
	return scores
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
