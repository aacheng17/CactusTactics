package main

type Message struct {
	client      *Client
	messageType int
	data        []byte
}

func newMessage(c *Client, t int, d []byte) *Message {
	return &Message{
		client:      c,
		messageType: t,
		data:        d,
	}
}
