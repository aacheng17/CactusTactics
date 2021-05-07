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
}

func (h *SpecializedHub) genNextLetters() {
	h.start = letters[rand.Intn(len(letters))]
	h.end = letters[rand.Intn(len(letters))]
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
	}
	h.genNextLetters()
	return h
}
