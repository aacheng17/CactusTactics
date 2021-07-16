package timeline

import (
	"fmt"
	"strconv"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type TimelineHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	playerNum int

	phase int

	events map[string]Event

	event Event
}

func (h *TimelineHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
	}
}

func (h *TimelineHub) getAssertedClients() map[*TimelineClient]bool {
	ret := make(map[*TimelineClient]bool)
	for k, v := range h.Clients {
		ret[k.(*TimelineClient)] = v
	}
	return ret
}

// RECEIVING:
// -: disconnect
// 0: name
// 1: lobby chat message
// 2: end game message
// a: guess

// SENDING:
// 0: restart
// 1: lobby chat message
// 2: end game
// 3: players
// a: prompt
// b: winners
// c: who's turn
// d: choice

func (h *TimelineHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*TimelineClient)
	if c.Name == "" && m.MessageType == byte('0') {
		name := m.Data[0]
		avatar, err1 := strconv.Atoi(m.Data[1])
		color, err2 := strconv.Atoi(m.Data[2])
		if err1 != nil || err2 != nil {
			return
		}
		c.Name = name
		c.Avatar = avatar
		c.Color = color
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
		h.SendData(c, byte('a'), h.getPrompt())
		h.SendData(c, byte('c'), h.getTurn())
		return
	}
	switch m.MessageType {
	case byte('1'):
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}
	if h.phase == -1 {
		if m.MessageType == byte('2') {
			h.reset()
			h.Broadcast(byte('0'), []string{""})
			h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " restarted the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(byte('3'), h.getPlayers())
			h.Broadcast(byte('a'), h.getPrompt())
			h.phase = 0
		}
		return
	}
	switch m.MessageType {
	case byte('a'):
		//WAKAWAKA
	case byte('2'):
		h.Broadcast(byte('2'), []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
		h.Broadcast(byte('b'), h.getWinners())
		h.phase = -1
	}
}

func NewTimelineHub() core.Hublike {
	h := &TimelineHub{
		Hub: *core.NewHub(),
	}
	h.Child = h
	h.reset()
	return h
}
