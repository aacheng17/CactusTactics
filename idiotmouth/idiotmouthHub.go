package idiotmouth

import (
	"fmt"
	"strings"

	"example.com/hello/core"
)

// declaring a struct
type IdiotmouthHub struct {

	// declaring struct variable
	core.Hub

	start rune

	end rune

	usedWords []string

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

func (h *IdiotmouthHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*IdiotmouthClient)
	switch m.MessageType {
	case byte('0'):
		for client := range h.Clients {
			h.SendData(client, byte('0'), []byte(c.Name+": "+string(m.Data)))
		}
	case byte('4'):
		word := string(m.Data)
		if definition, ok := dictionary[word]; ok {
			for client := range h.Clients {
				h.SendData(client, byte('0'), []byte(fmt.Sprint(c.Name+" asked, \"What?\" for the word ", word, ".<br/>", word, " - ", definition)))
			}
		}
	}
	if h.phase == 1 {
		if m.MessageType == byte('3') {
			h.reset()
			for client := range h.Clients {
				h.SendData(client, byte('4'), []byte(""))
				h.SendData(client, byte('0'), []byte(c.Name+" restarted the game<br/>"))
				h.SendData(client, byte('1'), []byte(h.getPlayers()))
				h.SendData(client, byte('2'), []byte(h.getPrompt()))
			}
			h.phase = 0
		}
		return
	}
	switch m.MessageType {
	case byte('0'):
		word := strings.TrimSpace(strings.ToLower(string(m.Data)))
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
					for client := range h.Clients {
						h.SendData(client, byte('0'), []byte("All possible words have been used or passed. Type restart to restart the game."))
					}
					break
				}
				for client := range h.Clients {
					h.SendData(client, byte('5'), []byte(fmt.Sprint(c.Name, " earned ", worth, "x", bonus, "=", finalWorth, " points for |", word)))
					h.SendData(client, byte('2'), []byte(h.getPrompt()))
					h.SendData(client, byte('1'), []byte(h.getPlayers()))
				}
			case 2:
				for client := range h.Clients {
					h.SendData(client, byte('0'), []byte("This word has already been used this game."))
				}
			}
		}
	case byte('1'):
		name := string(m.Data)
		if c.Name == "" {
			c.Name = name
		}
		for client := range h.Clients {
			h.SendData(client, byte('0'), []byte(name+" joined"))
		}
		for client := range h.Clients {
			h.SendData(client, byte('1'), []byte(h.getPlayers()))
		}
		h.SendData(c, byte('2'), []byte(h.getPrompt()))
	case byte('2'):
		if !c.pass {
			c.pass = true
			for client := range h.Clients {
				h.SendData(client, byte('0'), []byte(c.Name+" voted to skip."))
			}
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					for client := range h.Clients {
						h.SendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.Clients {
					h.SendData(client, byte('0'), []byte("Majority has voted to skip. New letters generated<br/>a"))
					h.SendData(client, byte('2'), []byte(h.getPrompt()))
				}
			}
		}
	case byte('3'):
		for client := range h.Clients {
			h.SendData(client, byte('0'), []byte("Game ended by "+c.Name))
			h.SendData(client, byte('3'), []byte(h.getWinners()))
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
