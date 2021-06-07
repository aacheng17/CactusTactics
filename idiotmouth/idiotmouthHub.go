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

func (h *IdiotmouthHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*IdiotmouthClient)
	if c.Name == "" && m.MessageType == byte('1') {
		name := string(m.Data[0])
		if c.Name == "" {
			c.Name = name
		}
		mNum := h.useMessageNum()
		for client := range h.Clients {
			h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		}
		for client := range h.Clients {
			h.SendData(client, byte('1'), h.getPlayers())
		}
		h.SendData(c, byte('2'), h.getPrompt())
		return
	}
	switch m.MessageType {
	case byte('0'):
		mNum := h.useMessageNum()
		for client := range h.Clients {
			h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
		}
	case byte('4'):
		clientMessageNum, err := strconv.Atoi(m.Data[0])
		if err != nil {
			break
		}
		if word := h.whattedWords[clientMessageNum]; word != "" {
			h.whattedWords[clientMessageNum] = ""
			if definition, ok := dictionary[word]; ok {
				mNum := h.useMessageNum()
				for client := range h.Clients {
					h.SendData(client, byte('6'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG+" said \"What?\" for the word ", word, u.ENDTAG, u.Tag("br"), u.Tag("p"), word, " - ", definition, u.ENDTAG), fmt.Sprint(clientMessageNum)})
				}
			}
		}
	}
	if h.phase == 1 {
		if m.MessageType == byte('3') {
			h.reset()
			mNum := h.useMessageNum()
			for client := range h.Clients {
				h.SendData(client, byte('4'), []string{""})
				h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p postbr", mNum), u.Tag("b")+c.Name+u.ENDTAG, " restarted the game", u.ENDTAG, u.ENDTAG)})
				h.SendData(client, byte('1'), h.getPlayers())
				h.SendData(client, byte('2'), h.getPrompt())
			}
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
					mNum := h.useMessageNum()
					for client := range h.Clients {
						h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					}
					break
				}
				mNum := h.useMessageNum()
				h.whattedWords[mNum] = word
				for client := range h.Clients {
					h.SendData(client, byte('5'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", word, u.ENDTAG), word})
					h.SendData(client, byte('2'), h.getPrompt())
					h.SendData(client, byte('1'), h.getPlayers())
				}
			case 2:
				mNum := h.useMessageNum()
				for client := range h.Clients {
					h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), "This word has already been used this game.", u.ENDTAG)})
				}
			}
		}
	case byte('2'):
		if !c.pass {
			c.pass = true
			mNum := h.useMessageNum()
			for client := range h.Clients {
				h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " voted to skip.", u.ENDTAG)})
			}
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					mNum := h.useMessageNum()
					for client := range h.Clients {
						h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p", mNum), "All possible words have been used or passed. Type restart to restart the game.", u.ENDTAG)})
					}
					break
				}
				mNum := h.useMessageNum()
				for client := range h.Clients {
					h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p postbr", mNum), "Majority has voted to skip. New letters generated", u.ENDTAG)})
					h.SendData(client, byte('2'), h.getPrompt())
				}
			}
		}
	case byte('3'):
		mNum := h.useMessageNum()
		for client := range h.Clients {
			h.SendData(client, byte('0'), []string{fmt.Sprint(u.TagId("p prebr postbr", mNum), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
			h.SendData(client, byte('3'), h.getWinners())
		}
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
