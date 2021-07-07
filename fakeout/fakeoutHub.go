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

	questions []int

	question int

	phase int

	answers []*FakeoutClient
}

func (h *FakeoutHub) DisconnectClientMessage(c core.Clientlike) {
	if c.GetName() != "" {
		h.Broadcast(byte('1'), []string{fmt.Sprint(u.TagId("p", 0), u.Tag("b")+c.GetName()+u.ENDTAG, " disconnected", u.ENDTAG)})
		h.Broadcast(byte('3'), h.getPlayers())
	}
}

func (h *FakeoutHub) getAssertedClients() map[*FakeoutClient]bool {
	ret := make(map[*FakeoutClient]bool)
	for k, v := range h.Clients {
		ret[k.(*FakeoutClient)] = v
	}
	return ret
}

// RECEIVING:
// -: disconnect
// 0: name
// 1: lobby chat message
// 2: end game message
// a: response
// b: choice
// c: request prompt

// SENDING:
// 0: restart
// 1: lobby chat message
// 2: end game
// 3: players
// a: prompt
// b: choice response
// c: choices
// d: choices response
// e: results
// f: winners

func (h *FakeoutHub) HandleHubMessage(m *core.Message) {
	c := (m.Client).(*FakeoutClient)
	question := questions.getQuestion(h.question)
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
		if h.phase == 1 {
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
				h.SendData(client, byte('c'), toSend)
			}
		}
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
		if h.phase == 0 {
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
					h.SendData(c, byte('b'), []string{"-1"}) //answer is too close to actual answer
				} else {
					h.SendData(c, byte('b'), []string{"0"})
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
							h.SendData(client, byte('c'), toSend)
						}
						h.phase = 1
					}
				}
			}
		}
	case byte('b'):
		if h.phase == 1 {
			playerChoice := strings.TrimSpace(string(m.Data[0]))
			choiceIndex, err := strconv.Atoi(playerChoice)
			if err == nil {
				if h.answers[choiceIndex] == c {
					h.SendData(c, byte('d'), []string{"-1"}) //can't pick your own number
				} else {
					h.SendData(c, byte('d'), []string{"0"})
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
								}
							}
						}
						results := []string{}
						for i, client := range h.answers {
							if client == nil {
								results = append(results, "ACTUAL ANSWER")
							} else {
								results = append(results, client.Name)
							}
							s := ""
							for _, fakedOut := range choices[i] {
								s += " " + fakedOut.Name
							}
							results = append(results, s)
						}
						h.phase = 0
						h.resetAnswers()
						h.genNextQuestion()
						for client := range h.Clients {
							h.SendData(client, byte('e'), results)
							h.SendData(client, byte('3'), h.getPlayers())
						}
						h.phase = 0
					}
				}
			}
		}
	case byte('c'):
		h.SendData(c, byte('a'), h.getPrompt())
		break
	case byte('2'):
		h.Broadcast(byte('2'), []string{fmt.Sprint(u.TagId("p prebr postbr", h.useMessageNum()), "Game ended by ", u.Tag("b")+c.Name+u.ENDTAG, u.ENDTAG)})
		h.Broadcast(byte('f'), h.getWinners())
		h.phase = -1
	}
}

func NewFakeoutHub() core.Hublike {
	h := &FakeoutHub{
		Hub: *core.NewHub(),
	}
	h.Child = h
	h.reset()
	return h
}
