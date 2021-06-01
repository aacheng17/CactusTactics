package core

type Message struct {
	Client      Clientlike
	MessageType byte
	Data        []byte
}

func NewMessage(c Clientlike, t byte, d []byte) *Message {
	return &Message{
		Client:      c,
		MessageType: t,
		Data:        d,
	}
}
