package main

import (
	"fmt"
	"math/rand"
)

// declaring a struct
type SpecializedHub struct {

	// declaring struct variable
	Hub

	start rune

	end rune

	usedWords []string
}

func (h *SpecializedHub) isWord(str string) bool {
	for _, v := range h.usedWords {
		if v == str {
			return false
		}
	}
	for _, v := range words {
		if v == str {
			return true
		}
	}
	return false
}

func (h *SpecializedHub) useWord(word string) {
	h.usedWords = append(h.usedWords, word)
}

func (h *SpecializedHub) reset() {
	for client := range h.clients {
		client.score = 0
		client.pass = false
	}
	h.genNextLetters()
}

func (h *SpecializedHub) resetPass() {
	for client := range h.clients {
		client.pass = false
	}
}

func (h *SpecializedHub) getMajorityPass() bool {
	count := 0
	for client := range h.clients {
		if client.pass {
			count++
		}
	}
	return count*2 >= len(h.clients)
}

func (h *SpecializedHub) genNextLetters() {
	h.start = letters[rand.Intn(len(letters))]
	h.end = letters[rand.Intn(len(letters))]
	h.start = 'b'
	h.end = 'o'
}

func (h *SpecializedHub) getWorth() int {
	return int(100.0 * (1.0 - float32(startFreq[h.start]+endFreq[h.start])/3000.0))
}

func (h *SpecializedHub) getPrompt() string {
	ret := string(h.start) + "*" + string(h.end)
	bonus := h.getWorth()
	ret += ", worth " + fmt.Sprint(bonus) + " points"
	return ret
}

func (h *SpecializedHub) getScores() string {
	scores := ""
	for client := range h.clients {
		scores += client.name + ": " + fmt.Sprint(client.score) + "; "
	}
	return scores
}

func newHub() *SpecializedHub {
	h := &SpecializedHub{
		Hub: Hub{
			register:   make(chan *SpecializedClient),
			unregister: make(chan *SpecializedClient),
			messages:   make(chan *Message),
			clients:    make(map[*SpecializedClient]bool),
		},
		usedWords: []string{},
	}
	h.genNextLetters()
	return h
}
