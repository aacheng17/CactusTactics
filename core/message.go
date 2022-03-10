package core

type Message struct {
	Client      Clientlike
	MessageCode byte
	Data        []string
}

func NewMessage(c Clientlike, t byte, d []string) *Message {
	return &Message{
		Client:      c,
		MessageCode: t,
		Data:        d,
	}
}
