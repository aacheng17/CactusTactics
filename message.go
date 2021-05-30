package main

type Message struct {
	client      Clientlike
	messageType byte
	data        []byte
}

func newMessage(c Clientlike, t byte, d []byte) *Message {
	return &Message{
		client:      c,
		messageType: t,
		data:        d,
	}
}
