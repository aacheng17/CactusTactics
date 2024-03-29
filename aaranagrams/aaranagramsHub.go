package aaranagrams

import (
	"fmt"
	"strconv"
	"time"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type AaranagramsHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	minWordLength int

	scoreToWin int

	chaosMode bool

	letters []rune

	phase byte

	turn int

	dictionary AaranagramsDictionary
}

func (h *AaranagramsHub) DisconnectClientMessage(c core.Clientlike) {
	if h.turn >= len(h.Clients) {
		h.turn = 0
		if len(h.Clients) == 0 {
			return
		}
	}
	h.sendYourTurn(true)
	if c.GetName() != "" {
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
	}
}

func (h *AaranagramsHub) getAssertedClients() map[*AaranagramsClient]bool {
	ret := make(map[*AaranagramsClient]bool)
	for k, v := range h.Clients {
		ret[k.(*AaranagramsClient)] = v
	}
	return ret
}

func (h *AaranagramsHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*AaranagramsClient)

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
		c.JoinTime = time.Now().UnixNano()
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		h.SendData(c, ToClientCode["IN_MEDIA_RES"], []string{string(h.phase)})
		h.SendData(c, ToClientCode["MIN_WORD_LENGTH"], []string{fmt.Sprint(h.minWordLength)})
		h.SendData(c, ToClientCode["SCORE_TO_WIN"], []string{fmt.Sprint(h.scoreToWin)})
		h.SendData(c, ToClientCode["CHAOS_MODE"], []string{h.getChaosModeAsString()})
		h.Broadcast(ToClientCode["LETTERS"], []string{string(h.letters)})
		if h.phase == Phase["PLAY"] {
			h.SendData(c, ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+"Minimum word length: "+u.ENDTAG, h.minWordLength, u.ENDTAG)})
			h.SendData(c, ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+"Score to win: "+u.ENDTAG, h.scoreToWin, u.ENDTAG)})
		}
	case ToServerCode["WHAT"]:
		clientMessageNum, err := strconv.Atoi(m.Data[0])
		if err != nil {
			break
		}
		if word := h.dictionary.whattedWords[clientMessageNum]; word != "" {
			h.dictionary.whattedWords[clientMessageNum] = ""
			if definition, ok := h.dictionary.dictionary[word]; ok {
				h.Broadcast(ToClientCode["WHAT_RESPONSE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG+" said \"What?\" for the word ", word, u.ENDTAG, u.Tag("p"), word, " - ", definition, u.ENDTAG), fmt.Sprint(clientMessageNum)})
			}
		}
	case ToServerCode["LOBBY_CHAT_MESSAGE"]:
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}

	switch h.phase {
	case Phase["PREGAME"]:
		switch m.MessageCode {
		case ToServerCode["MIN_WORD_LENGTH"]:
			minWordLength, err := strconv.Atoi(m.Data[0])
			if err != nil || minWordLength < 1 || minWordLength > 8 {
				break
			}
			h.minWordLength = minWordLength
			h.Broadcast(ToClientCode["MIN_WORD_LENGTH"], []string{fmt.Sprint(minWordLength)})
		case ToServerCode["SCORE_TO_WIN"]:
			scoreToWin, err := strconv.Atoi(m.Data[0])
			if err != nil || scoreToWin < 500 || scoreToWin > 50000 {
				break
			}
			h.scoreToWin = scoreToWin
			h.Broadcast(ToClientCode["SCORE_TO_WIN"], []string{fmt.Sprint(scoreToWin)})
		case ToServerCode["CHAOS_MODE"]:
			if m.Data[0] != "0" && m.Data[0] != "1" {
				break
			}
			originalChaosMode := h.chaosMode
			h.chaosMode = m.Data[0] == "1"
			if originalChaosMode != h.chaosMode {
				h.Broadcast(ToClientCode["CHAOS_MODE"], []string{h.getChaosModeAsString()})
			}
		case ToServerCode["START_GAME"]:
			h.reset()
			h.Broadcast(ToClientCode["START_GAME"], []string{""})
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " started the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(ToClientCode["LETTERS"], []string{string(h.letters)})
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+"Minimum word length: "+u.ENDTAG, h.minWordLength, u.ENDTAG)})
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+"Score to win: "+u.ENDTAG, h.scoreToWin, u.ENDTAG)})
			if h.chaosMode {
				h.broadcastYourTurn(true)
			} else {
				h.sendYourTurn(true)
			}
			h.phase = Phase["PLAY"]
			h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		}
	case Phase["PLAY"]:
		switch m.MessageCode {
		case ToServerCode["GAME_MESSAGE"]:
			if len(m.Data[0]) == 0 {
				break
			}
			h.handleWord(c, m.Data[0])
		case ToServerCode["LETTER"]:
			if h.chaosMode || h.getClientOfCurrentTurn() == c {
				letterCreated := false
				for i, l := range h.letters {
					if l == ' ' {
						h.letters[i] = u.GetLetterWeighted()
						h.Broadcast(ToClientCode["LETTERS"], []string{string(h.letters)})
						letterCreated = true
						break
					}
				}
				if !letterCreated {
					return
				}
				h.sendYourTurn(false)
				h.turn++
				if h.turn >= len(h.Clients) {
					h.turn = 0
				}
				h.sendYourTurn(true)
				h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
			}
		case ToServerCode["END_GAME"]:
			h.broadcastYourTurn(false)
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.Tag("p prebr"), u.Tag("b")+c.Name+u.ENDTAG, " ended the game.", u.ENDTAG)})
			h.endGame()
		}
	}
}

func NewAaranagramsHub(game string, id string, deleteHubCallback func(*core.Hub)) core.Hublike {
	h := &AaranagramsHub{
		Hub:           *core.NewHub(game, id, deleteHubCallback),
		phase:         Phase["PREGAME"],
		minWordLength: 3,
		scoreToWin:    3000,
		chaosMode:     false,
		letters:       make([]rune, 20),
		turn:          0,
	}
	h.Child = h
	return h
}
