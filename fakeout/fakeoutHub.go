package fakeout

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"example.com/hello/core"
	u "example.com/hello/utility"
)

// declaring a struct
type FakeoutHub struct {

	// declaring struct variable
	core.Hub

	messageNum int

	scoreToWin int

	deck int

	questions []int

	question int

	phase byte

	answers []*FakeoutClient
}

func (h *FakeoutHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", 0), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
	}
}

func (h *FakeoutHub) getAssertedClients() map[*FakeoutClient]bool {
	ret := make(map[*FakeoutClient]bool)
	for k, v := range h.Clients {
		ret[k.(*FakeoutClient)] = v
	}
	return ret
}

func (h *FakeoutHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*FakeoutClient)
	question := h.getQuestion()
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
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+name+u.ENDTAG, " joined", u.ENDTAG)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		h.SendData(c, ToClientCode["IN_MEDIA_RES"], []string{string(h.phase)})
		switch h.phase {
		case Phase["PREGAME"]:
			h.SendData(c, ToClientCode["SCORE_TO_WIN"], []string{fmt.Sprint(h.scoreToWin)})
			h.SendData(c, ToClientCode["DECK_OPTIONS"], getFakeoutDeckOptions())
			h.SendData(c, ToClientCode["DECK_SELECTION"], []string{fmt.Sprint(h.deck)})
		case Phase["PLAY_PROMPT"]:
			h.SendData(c, ToClientCode["PROMPT"], h.getPrompt())
		case Phase["PLAY_GUESSES"]:
			h.SendData(c, ToClientCode["PROMPT"], h.getPrompt())
			toSend := []string{}
			for _, client := range h.answers {
				s := ""
				if client == nil {
					s = question.Answer
				} else {
					s = client.answer
				}
				toSend = append(toSend, s)
			}
			h.SendData(c, ToClientCode["CHOICES"], toSend)
		}
	case ToServerCode["LOBBY_CHAT_MESSAGE"]:
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, ": ", m.Data[0], u.ENDTAG)})
	}

	switch h.phase {
	case Phase["PREGAME"]:
		switch m.MessageCode {
		case ToServerCode["SCORE_TO_WIN"]:
			scoreToWin, err := strconv.Atoi(m.Data[0])
			if err != nil || scoreToWin < 250 || scoreToWin > 50000 {
				break
			}
			h.scoreToWin = scoreToWin
			fmt.Println("wope")
			fmt.Println(scoreToWin)
			h.Broadcast(ToClientCode["SCORE_TO_WIN"], []string{fmt.Sprint(scoreToWin)})
		case ToServerCode["DECK_SELECTION"]:
			deckSelection, err := strconv.Atoi(m.Data[0])
			if err == nil {
				h.deck = deckSelection
				h.Broadcast(ToClientCode["DECK_SELECTION"], []string{fmt.Sprint(h.deck)})
			}
		case ToServerCode["START_GAME"]:
			h.startGame()
			h.Broadcast(ToClientCode["START_GAME"], []string{""})
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p postbr", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " started the game", u.ENDTAG, u.ENDTAG)})
			h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
			h.phase = Phase["PLAY_PROMPT"]
		}
	case Phase["PLAY_PROMPT"]:
		switch m.MessageCode {
		case ToServerCode["RESPONSE"]:
			if c.answer == "" {
				playerAnswer := strings.TrimSpace(strings.ToLower(string(m.Data[0])))
				alternateSpelling := false
				for _, x := range question.AlternateSpellings {
					if playerAnswer == x {
						alternateSpelling = true
						break
					}
				}
				if c.answer != "" {
				} else if playerAnswer == question.Answer || alternateSpelling {
					h.SendData(c, ToClientCode["CHOICE_RESPONSE"], []string{"-1"}) //answer is too close to actual answer
				} else {
					h.SendData(c, ToClientCode["CHOICE_RESPONSE"], []string{"0"})
					c.answer = playerAnswer
					if h.isAllAnswered() {
						h.answers = []*FakeoutClient{nil}
						for client := range h.getAssertedClients() {
							h.answers = append(h.answers, client)
						}
						rand.Shuffle(len(h.answers), func(i, j int) { h.answers[i], h.answers[j] = h.answers[j], h.answers[i] })
						toSend := []string{}
						for _, client := range h.answers {
							s := ""
							if client == nil {
								s = question.Answer
							} else {
								s = client.answer
							}
							toSend = append(toSend, s)
						}
						for client := range h.Clients {
							h.SendData(client, ToClientCode["CHOICES"], toSend)
						}
						h.phase = Phase["PLAY_GUESSES"]
					}
					h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
				}
			}
		}
	case Phase["PLAY_GUESSES"]:
		switch m.MessageCode {
		case ToServerCode["CHOICE"]:
			playerChoice := strings.TrimSpace(string(m.Data[0]))
			choiceIndex, err := strconv.Atoi(playerChoice)
			if err == nil {
				if h.answers[choiceIndex] == c {
					h.SendData(c, ToClientCode["CHOICES_RESPONSE"], []string{"-1"}) //can't pick your own number
				} else {
					h.SendData(c, ToClientCode["CHOICES_RESPONSE"], []string{"0"})
					c.choice = choiceIndex
					if h.isAllChosen() {
						choices := make([][]*FakeoutClient, len(h.answers))
						for i := range choices {
							choices[i] = make([]*FakeoutClient, 0)
						}
						for client := range h.getAssertedClients() {
							choices[client.choice] = append(choices[client.choice], client)
							if h.answers[client.choice] == nil {
								client.score += 100
							} else {
								if h.answers[client.choice] != client {
									h.answers[client.choice].score += 50
									h.answers[client.choice].fakeouts++
								}
							}
						}
						results := []string{}
						for i, client := range h.answers {
							if client == nil {
								results = append(results, "")
							} else {
								results = append(results, client.Name)
							}
							s := ""
							for _, fakedOut := range choices[i] {
								s += " " + fakedOut.Name
							}
							results = append(results, s)
						}
						h.resetAnswers()
						h.Broadcast(ToClientCode["RESULTS"], results)
						h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
						if h.didSomeoneWin() {
							h.Broadcast(ToClientCode["WINNERS"], h.getWinners())
							h.phase = Phase["PREGAME"]
						} else {
							h.phase = Phase["PLAY_PROMPT"]
							h.genNextQuestion()
						}
					}
					h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
				}
			}
		}
	}
	if h.phase != Phase["PREGAME"] {
		switch m.MessageCode {
		case ToServerCode["PROMPT_REQUEST"]:
			h.SendData(c, ToClientCode["PROMPT"], h.getPrompt())
		case ToServerCode["END_GAME"]:
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
			h.Broadcast(ToClientCode["WINNERS"], h.getWinners())
			h.phase = Phase["PREGAME"]
		}
	}
}

func NewFakeoutHub(game string, id string, deleteHubCallback func(*core.Hub)) core.Hublike {
	h := &FakeoutHub{
		Hub:        *core.NewHub(game, id, deleteHubCallback),
		scoreToWin: 1000,
	}
	h.Child = h
	h.reset()
	return h
}
