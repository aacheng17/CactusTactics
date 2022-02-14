package idiotmouth

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type IdiotmouthHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	start rune

	end rune

	usedWords map[string]bool

	whattedWords map[int]string

	letters map[string]int

	wordsLeft int

	phase byte
}

func (h *IdiotmouthHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
	}
}

func (h *IdiotmouthHub) getAssertedClients() map[*IdiotmouthClient]bool {
	ret := make(map[*IdiotmouthClient]bool)
	for k, v := range h.Clients {
		ret[k.(*IdiotmouthClient)] = v
	}
	return ret
}

func (h *IdiotmouthHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*IdiotmouthClient)
	if c.Name == "" && m.MessageType == ToServerCode["NAME"] {
		name := m.Data[0]
		avatar, err1 := strconv.Atoi(m.Data[1])
		color, err2 := strconv.Atoi(m.Data[2])
		if err1 != nil || err2 != nil {
			return
		}
		c.Name = name
		c.Avatar = avatar
		c.Color = color
		h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		h.SendData(c, ToClientCode["PROMPT"], h.getPrompt())
		if h.phase == Phase["PREGAME"] {
			h.SendData(c, ToClientCode["END_GAME"], []string{})
		}
		return
	}
	switch m.MessageType {
	case ToServerCode["WHAT"]:
		clientMessageNum, err := strconv.Atoi(m.Data[0])
		if err != nil {
			break
		}
		if word := h.whattedWords[clientMessageNum]; word != "" {
			h.whattedWords[clientMessageNum] = ""
			if definition, ok := dictionary[word]; ok {
				h.Broadcast(ToClientCode["WHAT_RESPONSE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG+" said \"What?\" for the word ", word, u.ENDTAG, u.Tag("p"), word, " - ", definition, u.ENDTAG), fmt.Sprint(clientMessageNum)})
			}
		}
	case ToServerCode["LOBBY_CHAT_MESSAGE"]:
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}
	if h.phase == Phase["PREGAME"] {
		if m.MessageType == ToServerCode["END_GAME"] {
			h.reset()
			h.Broadcast(ToClientCode["RESTART"], []string{""})
			h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " restarted the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
			h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
			h.phase = Phase["PLAY"]
		}
		return
	}
	switch m.MessageType {
	case ToServerCode["GAME_MESSAGE"]:
		word := strings.TrimSpace(strings.ToLower(string(m.Data[0])))
		h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", word)})
		if len(word) >= 3 && word[0] == byte(h.start) && word[len(word)-1] == byte(h.end) {
			switch h.validWord(word) {
			case 0:
				worth := h.getWorth()
				bonus := len(word) - 2
				finalWorth := worth * bonus
				c.score += finalWorth
				if finalWorth > c.highestScore {
					c.highestWord = word
					c.highestScore = finalWorth
				}
				err := h.gotIt(word)
				if err == 1 {
					h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				mNum := h.useMessageNum()
				h.whattedWords[mNum] = word
				h.Broadcast(ToClientCode["MESSAGE_WITH_WHAT"], []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", word, u.ENDTAG), word})
				h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
				h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
			case 2:
				h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "This word has already been used this game.", u.ENDTAG)})
			}
		}
	case ToServerCode["VOTE_SKIP"]:
		if !c.pass {
			c.pass = true
			h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " voted to skip.", u.ENDTAG)})
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), "Majority has voted to skip. New letters generated", u.ENDTAG)})
				h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
			}
		}
	case ToServerCode["END_GAME"]:
		h.Broadcast(ToClientCode["END_GAME"], []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
		h.Broadcast(ToClientCode["WINNERS"], h.getWinners())
		h.phase = Phase["PREGAME"]
	}
}

func NewIdiotmouthHub(game string, id string, deleteHubCallback func(*core.Hub)) core.Hublike {
	h := &IdiotmouthHub{
		Hub:     *core.NewHub(game, id, deleteHubCallback),
		letters: make(map[string]int),
	}
	h.Child = h
	h.reset()
	return h
}
