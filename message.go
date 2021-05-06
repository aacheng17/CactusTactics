package main

type Message struct {
	client      *Client
	messageType byte
	data        []byte
}

func newMessage(c *Client, t byte, d []byte) *Message {
	return &Message{
		client:      c,
		messageType: t,
		data:        d,
	}
}
