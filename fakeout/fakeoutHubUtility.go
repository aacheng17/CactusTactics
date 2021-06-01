package fakeout

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"example.com/hello/utility"
)

func (h *FakeoutHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.answer = ""
		client.choice = -1
	}
	h.phase = 0
	h.questions = utility.MakeRange(0, questions.size())
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
	for i, x := range h.questions {
		if x == h.question {
			h.questions[i] = h.questions[len(h.questions)-1]
			h.questions = h.questions[:len(h.questions)-1]
		}
	}
	h.question = h.questions[rand.Intn(len(h.questions))]
	return 0
}

func (h *FakeoutHub) getPrompt() string {
	ret := questions.getQuestion(h.question).Question
	ret = strings.Replace(ret, "<BLANK>", "________", 1)
	return ret
}

func (h *FakeoutHub) getScores() string {
	keys := make([]*FakeoutClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	scores := ""
	for _, client := range keys {
		if client.Name == "" {
			continue
		}
		scores += client.Name + ": " + fmt.Sprint(client.score) + "; "
	}
	return scores
}
