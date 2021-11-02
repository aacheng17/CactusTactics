package standoff

import (
	"fmt"
	"strconv"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type StandoffHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	phase int

	nextClientId int

	round int
}

func (h *StandoffHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
	}
}

func (h *StandoffHub) getAssertedClients() map[*StandoffClient]bool {
	ret := make(map[*StandoffClient]bool)
	for k, v := range h.Clients {
		ret[k.(*StandoffClient)] = v
	}
	return ret
}

// RECEIVING:
// -: disconnect
// 0: name
// 1: lobby chat message
// 2: end game message
// a: game decision
// b: prompt request

// SENDING:
// 0: restart
// 1: lobby chat message
// 2: end game
// 3: players
// a: round prompt and choices
// b: choice ack
// c: result
// d: winners

func (h *StandoffHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*StandoffClient)
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
		c.id = h.nextClientId
		h.nextClientId++
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
		if h.phase == -1 {
			h.SendData(c, byte('2'), []string{""})
		}
		return
	}
	switch m.MessageType {
	case byte('1'):
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}
	if h.phase == -1 {
		switch m.MessageType {
		case byte('2'):
			h.reset()
			h.Broadcast(byte('0'), []string{""})
			h.Broadcast(byte('3'), h.getPlayers())
			h.Broadcast(byte('a'), h.getPrompt())
		}
	} else if h.phase == 0 {
		switch m.MessageType {
		case byte('a'):
			decision, err := strconv.Atoi(string(m.Data[0]))
			if err != nil {
				return
			}
			h.SendData(c, 'b', []string{""})
			if c.decision == -1 {
				c.decision = decision
				h.calcDecisionResult()
			}
		case byte('b'):
			if !c.active {
				h.SendData(c, byte('a'), []string{fmt.Sprint(h.round), "spectating"})
			} else if !c.alive {
				h.SendData(c, byte('a'), []string{fmt.Sprint(h.round), "dead"})
			} else {
				h.SendData(c, byte('a'), h.getPrompt())
			}
		}
	}
}

func (h *StandoffHub) calcDecisionResult() {
	if h.isAllDecided() {
		h.Broadcast(byte('c'), h.calcResult())
		h.Broadcast(byte('3'), h.getPlayers())
		if h.numAlive() < 2 {
			h.Broadcast(byte('2'), []string{""})
			h.Broadcast(byte('d'), h.getWinners())
			h.phase = -1
		} else {
			h.nextRound()
		}
	}
}

func NewStandoffHub(game string, id string, deleteHubCallback func(*core.Hub)) core.Hublike {
	h := &StandoffHub{
		Hub:          *core.NewHub(game, id, deleteHubCallback),
		phase:        -1,
		round:        0,
		nextClientId: 0,
	}
	h.Child = h
	h.reset()
	return h
}
