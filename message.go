package main

type Message struct {
	client      *SpecializedClient
	messageType byte
	data        []byte
}

func newMessage(c *SpecializedClient, t byte, d []byte) *Message {
	return &Message{
		client:      c,
		messageType: t,
		data:        d,
	}
}
