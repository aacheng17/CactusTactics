package main

import "math/rand"

// declaring a struct
type SpecializedHub struct {

	// declaring struct variable
	Hub

	start rune

	end rune
}

func newHub() *SpecializedHub {
	return &SpecializedHub{
		Hub: Hub{
			register:   make(chan *SpecializedClient),
			unregister: make(chan *SpecializedClient),
			messages:   make(chan *Message),
			clients:    make(map[*SpecializedClient]bool),
		},
		start: letters[rand.Intn(len(letters))],
		end:   letters[rand.Intn(len(letters))],
	}
}
