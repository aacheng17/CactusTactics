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

	phase int
}

func (h *IdiotmouthHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(byte('1'), h.getPlayers())
	}
}

func (h *IdiotmouthHub) getAssertedClients() map[*IdiotmouthClient]bool {
	ret := make(map[*IdiotmouthClient]bool)
	for k, v := range h.Clients {
		ret[k.(*IdiotmouthClient)] = v
	}
	return ret
}

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: winners
// 4: restart (data is inconsequential, probably empty string)
// 5: message that needs a "what?""
// 6: "what?" message
// 7: end game message

func (h *IdiotmouthHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*IdiotmouthClient)
	if c.Name == "" && m.MessageType == byte('1') {
		name := m.Data[0]
		avatar, err1 := strconv.Atoi(m.Data[1])
		color, err2 := strconv.Atoi(m.Data[2])
		if err1 != nil || err2 != nil {
			return
		}
		c.Name = name
		c.Avatar = avatar
		c.Color = color
		h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(byte('1'), h.getPlayers())
		h.SendData(c, byte('2'), h.getPrompt())
		return
	}
	switch m.MessageType {
	case byte('0'):
		h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	case byte('4'):
		clientMessageNum, err := strconv.Atoi(m.Data[0])
		if err != nil {
			break
		}
		if word := h.whattedWords[clientMessageNum]; word != "" {
			h.whattedWords[clientMessageNum] = ""
			if definition, ok := dictionary[word]; ok {
				h.Broadcast(byte('6'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG+" said \"What?\" for the word ", word, u.ENDTAG, u.Tag("p"), word, " - ", definition, u.ENDTAG), fmt.Sprint(clientMessageNum)})
			}
		}
	case byte('8'):
		h.Broadcast(byte('8'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}
	if h.phase == 1 {
		if m.MessageType == byte('3') {
			h.reset()
			h.Broadcast(byte('4'), []string{""})
			h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " restarted the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(byte('1'), h.getPlayers())
			h.Broadcast(byte('2'), h.getPrompt())
			h.phase = 0
		}
		return
	}
	switch m.MessageType {
	case byte('0'):
		word := strings.TrimSpace(strings.ToLower(string(m.Data[0])))
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
					h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				mNum := h.useMessageNum()
				h.whattedWords[mNum] = word
				h.Broadcast(byte('5'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", word, u.ENDTAG), word})
				h.Broadcast(byte('2'), h.getPrompt())
				h.Broadcast(byte('1'), h.getPlayers())
			case 2:
				h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "This word has already been used this game.", u.ENDTAG)})
			}
		}
	case byte('2'):
		if !c.pass {
			c.pass = true
			h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " voted to skip.", u.ENDTAG)})
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				h.Broadcast(byte('0'), []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), "Majority has voted to skip. New letters generated", u.ENDTAG)})
				h.Broadcast(byte('2'), h.getPrompt())
			}
		}
	case byte('3'):
		h.Broadcast(byte('7'), []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
		h.Broadcast(byte('3'), h.getWinners())
		h.phase = 1
	}
}

func NewIdiotmouthHub() core.Hublike {
	h := &IdiotmouthHub{
		Hub:     *core.NewHub(),
		letters: make(map[string]int),
	}
	h.Child = h
	h.reset()
	return h
}
