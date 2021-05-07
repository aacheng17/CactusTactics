package main

// declaring a struct
type SpecializedHub struct {

	// declaring struct variable
	Hub

	score int
}

func newHub() *SpecializedHub {
	return &SpecializedHub{
		Hub: Hub{
			register:   make(chan *SpecializedClient),
			unregister: make(chan *SpecializedClient),
			messages:   make(chan *Message),
			clients:    make(map[*SpecializedClient]bool),
		},
	}
}
