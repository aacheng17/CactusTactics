package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

// declaring a struct
type FakeoutHub struct {

	// declaring struct variable
	Hub

	questions []int

	question int

	phase int

	answers []*FakeoutClient
}

func (h *FakeoutHub) getAssertedClients() map[*FakeoutClient]bool {
	ret := make(map[*FakeoutClient]bool)
	for k, v := range h.clients {
		ret[k.(*FakeoutClient)] = v
	}
	return ret
}

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: restart (data is inconsequential, probably empty string)

func (h *FakeoutHub) handleHubMessage(m *Message) {
	c := (m.client).(*FakeoutClient)
	question := questions.getQuestion(h.question)
	switch m.messageType {
	case byte('0'):
		if string(m.data) == "/restart" {
			h.reset()
			for client := range h.clients {
				h.sendData(client, byte('3'), []byte(""))
				h.sendData(client, byte('0'), []byte(c.name+" restarted the game"))
				h.sendData(client, byte('1'), []byte(h.getScores()))
				h.sendData(client, byte('2'), []byte(h.getPrompt()))
				h.sendData(client, byte('0'), []byte("."))
				h.sendData(client, byte('0'), []byte("New Prompt: "+h.getPrompt()))
			}
		} else if h.phase == 0 {
			playerAnswer := strings.TrimSpace(strings.ToLower(string(m.data)))
			alternateSpelling := false
			for _, x := range question.AlternateSpellings {
				if playerAnswer == x {
					alternateSpelling = true
					break
				}
			}
			if c.answer != "" {
			} else if playerAnswer == question.Answer || alternateSpelling {
				h.sendData(c, byte('0'), []byte("Your answer is too close to the actual answer. Please choose another answer."))
			} else {
				h.sendData(c, byte('0'), []byte("Your answer has been recorded. Waiting for other players' answers."))
				c.answer = playerAnswer
				if h.isAllAnswered() {
					h.answers = []*FakeoutClient{nil}
					for client := range h.getAssertedClients() {
						h.answers = append(h.answers, client)
					}
					rand.Shuffle(len(h.answers), func(i, j int) { h.answers[i], h.answers[j] = h.answers[j], h.answers[i] })
					stringToSend := "Choose from these answers:<br/>"
					for i, client := range h.answers {
						s := ""
						if client == nil {
							s = question.Answer
						} else {
							s = client.answer
						}
						stringToSend += "(" + fmt.Sprint(i) + ") " + s + "<br/>"
					}
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte(stringToSend))
					}
					h.phase = 1
				}
			}
		} else if h.phase == 1 {
			playerChoice := strings.TrimSpace(string(m.data))
			choiceIndex, err := strconv.Atoi(playerChoice)
			if err != nil || choiceIndex < 0 || choiceIndex >= len(h.answers) {
				h.sendData(c, byte('0'), []byte("Invalid choice. Please enter a valid number choice."))
			} else if h.answers[choiceIndex] == c {
				h.sendData(c, byte('0'), []byte("Invalid choice. You can't pick your own answer."))
			} else {
				c.choice = choiceIndex
				wordChoice := ""
				if h.answers[c.choice] == nil {
					wordChoice = question.Answer
				} else {
					wordChoice = h.answers[c.choice].answer
				}
				h.sendData(c, byte('0'), []byte("You chose ("+fmt.Sprint(c.choice)+") "+wordChoice+". Waiting for other players' choices."))
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
					stringToSend := "Results:<br/>"
					for i, client := range h.answers {
						if client == nil {
							stringToSend += question.Answer + " (ACTUAL ANSWER)"
						} else {
							stringToSend += client.answer + " (" + client.name + ") faked out"
							if len(choices[i]) == 0 {
								stringToSend += " no one"
							}
						}
						for _, fakedOut := range choices[i] {
							stringToSend += " " + fakedOut.name
						}
						stringToSend += "<br/>"
					}
					h.phase = 0
					h.resetAnswers()
					h.genNextQuestion()
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte(stringToSend))
						h.sendData(client, byte('1'), []byte(h.getScores()))
						h.sendData(client, byte('2'), []byte(h.getPrompt()))
						h.sendData(client, byte('0'), []byte("."))
						h.sendData(client, byte('0'), []byte("New Prompt: "+h.getPrompt()))
					}
				}
			}
		}
	case byte('1'):
		name := string(m.data)
		if c.name == "" {
			c.name = name
		}
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(name+" joined"))
		}
		for client := range h.clients {
			h.sendData(client, byte('1'), []byte(h.getScores()))
		}
		h.sendData(c, byte('2'), []byte(h.getPrompt()))
		h.sendData(c, byte('0'), []byte("New Prompt: "+h.getPrompt()))
	}
}

func (h *FakeoutHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.removeClient(client, "Removed client that disconnected.")
			}
		case message := <-h.messages:
			log.Println("Received message\n\tType: " + fmt.Sprint(message.messageType) + "\n\tData: " + string(message.data))
			h.handleHubMessage(message)
		}
	}
}

func newFakeoutHub() Hublike {
	h := &FakeoutHub{
		Hub: *newHub(),
	}
	h.reset()
	return h
}
