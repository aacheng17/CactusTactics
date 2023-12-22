package standoff

import (
	"fmt"
	"strconv"
	"time"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type StandoffHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	phase byte

	nextClientId int

	round int
}

func (h *StandoffHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
	}
}

func (h *StandoffHub) getAssertedClients() map[*StandoffClient]bool {
	ret := make(map[*StandoffClient]bool)
	for k, v := range h.Clients {
		ret[k.(*StandoffClient)] = v
	}
	return ret
}

func (h *StandoffHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*StandoffClient)
	switch m.MessageCode {
	case ToServerCode["NAME"]:
		if c.Name != "" {
			return
		}
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
		c.JoinTime = time.Now().UnixNano()
		h.nextClientId++
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		h.SendData(c, ToClientCode["IN_MEDIA_RES"], []string{string(h.phase)})
		return
	case ToServerCode["LOBBY_CHAT_MESSAGE"]:
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}

	switch h.phase {
	case Phase["PREGAME"]:
		switch m.MessageCode {
		case ToServerCode["START_GAME"]:
			h.reset()
			h.phase = Phase["PLAY"]
			h.Broadcast(ToClientCode["START_GAME"], []string{""})
			h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
		}
	case Phase["PLAY"]:
		switch m.MessageCode {
		case ToServerCode["DECISION"]:
			decision, err := strconv.Atoi(string(m.Data[0]))
			if err != nil {
				return
			}
			h.SendData(c, ToClientCode["DECISION_ACK"], []string{""})
			if c.decision == -1 {
				c.decision = decision
				h.calcDecisionResult()
			}
			h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		case ToServerCode["PROMPT_REQUEST"]:
			if !c.active {
				h.SendData(c, ToClientCode["PROMPT"], []string{fmt.Sprint(h.round), "spectating"})
			} else if !c.alive {
				h.SendData(c, ToClientCode["PROMPT"], []string{fmt.Sprint(h.round), "dead"})
			} else {
				h.SendData(c, ToClientCode["PROMPT"], h.getPrompt())
			}
		case ToServerCode["END_GAME"]:
			h.endGame()
		}
	}
}

func (h *StandoffHub) calcDecisionResult() {
	if h.isAllDecided() {
		h.Broadcast(ToClientCode["RESULT"], h.calcResult())
		for client := range h.getAssertedClients() {
			if client.alive {
				client.roundsAlive++
			}
		}
		if h.numAlive() < 2 {
			h.Broadcast(ToClientCode["END_GAME"], []string{""})
			h.Broadcast(ToClientCode["WINNERS"], h.getWinners())
			h.phase = Phase["PREGAME"]
		} else {
			h.nextRound()
		}
	}
}

func NewStandoffHub(game string, id string, deleteHubCallback func(*core.Hub)) core.Hublike {
	h := &StandoffHub{
		Hub:   *core.NewHub(game, id, deleteHubCallback),
		phase: Phase["PREGAME"],
	}
	h.Child = h
	h.reset()
	return h
}
