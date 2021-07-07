package timeline

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type TimelineHub struct {

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

func (h *TimelineHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
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
// a: standard game log message
// b: vote to skip
// c: what?

// SENDING:
// 0: restart
// 1: lobby chat message
// 2: end game
// 3: players
// a: standard game log message
// b: what response
// c: message containing a what
// d: prompt
// e: winners

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
		h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
		h.SendData(c, byte('d'), h.getPrompt())
		return
	}
	switch m.MessageType {
	case byte('c'):
		clientMessageNum, err := strconv.Atoi(m.Data[0])
		if err != nil {
			break
		}
		if word := h.whattedWords[clientMessageNum]; word != "" {
			h.whattedWords[clientMessageNum] = ""
			if definition, ok := dictionary[word]; ok {
				h.Broadcast(byte('b'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG+" said \"What?\" for the word ", word, u.ENDTAG, u.Tag("p"), word, " - ", definition, u.ENDTAG), fmt.Sprint(clientMessageNum)})
			}
		}
	case byte('1'):
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}
	if h.phase == -1 {
		if m.MessageType == byte('2') {
			h.reset()
			h.Broadcast(byte('0'), []string{""})
			h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " restarted the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(byte('3'), h.getPlayers())
			h.Broadcast(byte('d'), h.getPrompt())
			h.phase = 0
		}
		return
	}
	switch m.MessageType {
	case byte('a'):
		word := strings.TrimSpace(strings.ToLower(string(m.Data[0])))
		h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", word)})
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
					h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				mNum := h.useMessageNum()
				h.whattedWords[mNum] = word
				h.Broadcast(byte('c'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", word, u.ENDTAG), word})
				h.Broadcast(byte('d'), h.getPrompt())
				h.Broadcast(byte('3'), h.getPlayers())
			case 2:
				h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "This word has already been used this game.", u.ENDTAG)})
			}
		}
	case byte('b'):
		if !c.pass {
			c.pass = true
			h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " voted to skip.", u.ENDTAG)})
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					break
				}
				h.Broadcast(byte('a'), []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), "Majority has voted to skip. New letters generated", u.ENDTAG)})
				h.Broadcast(byte('d'), h.getPrompt())
			}
		}
	case byte('2'):
		h.Broadcast(byte('2'), []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
		h.Broadcast(byte('e'), h.getWinners())
		h.phase = -1
	}
}

func NewTimelineHub() core.Hublike {
	h := &TimelineHub{
		Hub:     *core.NewHub(),
		letters: make(map[string]int),
	}
	h.Child = h
	h.reset()
	return h
}
