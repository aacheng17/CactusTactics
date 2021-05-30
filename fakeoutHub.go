package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
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

func (h *FakeoutHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.answer = ""
		client.choice = -1
	}
	h.phase = 0
	h.questions = makeRange(0, questions.size())
	h.genNextQuestion()
}

func (h *FakeoutHub) resetAnswers() {
	for client := range h.getAssertedClients() {
		client.answer = ""
		client.choice = -1
	}
}

func (h *FakeoutHub) isAllAnswered() bool {
	for client := range h.getAssertedClients() {
		if client.answer == "" {
			return false
		}
	}
	return true
}

func (h *FakeoutHub) isAllChosen() bool {
	for client := range h.getAssertedClients() {
		if client.choice == -1 {
			return false
		}
	}
	return true
}

func (h *FakeoutHub) genNextQuestion() int {
	if len(h.questions) <= 0 {
		return 1
	}
	h.question = rand.Intn(len(h.questions))
	return 0
}

func (h *FakeoutHub) getPrompt() string {
	ret := questions.getQuestion(h.question).Question
	ret = strings.Replace(ret, "<BLANK>", "________", 1)
	return ret
}

func (h *FakeoutHub) getScores() string {
	keys := make([]*FakeoutClient, 0, len(h.clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	scores := ""
	for _, client := range keys {
		if client.name == "" {
			continue
		}
		scores += client.name + ": " + fmt.Sprint(client.score) + "; "
	}
	return scores
}

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: restart (data is inconsequential, probably empty string)

func (h *FakeoutHub) handleHubMessage(m *Message) {
	c := (m.client).(*FakeoutClient)
	switch m.messageType {
	case byte('0'):
		if string(m.data) == "restart" {
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
			question := questions.getQuestion(h.question)
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
			} else {
				c.choice = choiceIndex
				h.sendData(c, byte('0'), []byte("Your choice has been recorded. Waiting for other players' choices."))
				if h.isAllChosen() {
					question := questions.getQuestion(h.question)
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

func newFakeoutHub() *FakeoutHub {
	h := &FakeoutHub{
		Hub: *newHub(),
	}
	h.reset()
	return h
}
