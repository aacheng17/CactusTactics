package core

type Message struct {
	Client      Clientlike
	MessageType byte
	Data        []string
}

func NewMessage(c Clientlike, t byte, d []string) *Message {
	return &Message{
		Client:      c,
		MessageType: t,
		Data:        d,
	}
}
