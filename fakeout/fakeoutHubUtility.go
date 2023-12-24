package fakeout

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"example.com/hello/utility"
	u "example.com/hello/utility"
)

func (h *FakeoutHub) getQuestion() Question {
	return getFakeoutQuestion(h.deck, h.question)
}

func (h *FakeoutHub) useMessageNum() int {
	ret := h.messageNum
	h.messageNum++
	return ret
}

func (h *FakeoutHub) endGame() {
	for client := range h.getAssertedClients() {
		client.answer = ""
	}
	h.phase = Phase["PREGAME"]
}

func (h *FakeoutHub) reset() {
	for client := range h.getAssertedClients() {
		client.fakeouts = 0
		client.score = 0
		client.choice = -1
	}
}

func (h *FakeoutHub) startGame() {
	h.questions = utility.MakeRange(0, len(decks[h.deck].Questions))
	h.phase = Phase["PLAY_PROMPT"]
	h.question = h.questions[rand.Intn(len(h.questions))]
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

func (h *FakeoutHub) getPrompt() []string {
	ret := fmt.Sprint(decks[h.deck].Instructions, ":", u.Tag("br"), u.Tag("br"))
	ret += h.getQuestion().Question
	ret = strings.Replace(ret, "<BLANK>", "________", 1)
	ret = utility.ParseAndTag(ret)
	return []string{ret}
}

func (h *FakeoutHub) getPlayers(excepts ...*FakeoutClient) []string {
	keys := make([]*FakeoutClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		isExcept := false
		for _, e := range excepts {
			if k == e {
				isExcept = true
				break
			}
		}
		if !isExcept {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	players := []string{}
	for _, client := range keys {
		if client.Name == "" {
			continue
		}
		players = append(players, client.Name)
		players = append(players, fmt.Sprint(client.Avatar))
		players = append(players, fmt.Sprint((client.Color)))
		players = append(players, fmt.Sprint(client.score))
		players = append(players, fmt.Sprint(client.fakeouts))
		dotdotdotStatus := "none"
		if h.phase == Phase["PLAY_PROMPT"] {
			if client.answer == "" {
				dotdotdotStatus = "dotdotdot"
			} else {
				dotdotdotStatus = "ready"
			}
		} else if h.phase == Phase["PLAY_GUESSES"] {
			if client.choice == -1 {
				dotdotdotStatus = "dotdotdot"
			} else {
				dotdotdotStatus = "ready"
			}
		}
		players = append(players, fmt.Sprint(dotdotdotStatus))
	}
	return players
}

func (h *FakeoutHub) didSomeoneWin() bool {
	for k := range h.getAssertedClients() {
		if k.score >= h.scoreToWin {
			fmt.Println("yup")
			return true
		}
	}
	fmt.Println("nope")
	return false
}

func (h *FakeoutHub) getWinners() []string {
	ret := []string{}
	keys := make([]*FakeoutClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	winner := keys[0]
	ret = append(ret, winner.Name)
	ret = append(ret, fmt.Sprint(winner.score))

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].fakeouts > keys[j].fakeouts
	})
	winner = keys[0]
	ret = append(ret, winner.Name)
	ret = append(ret, fmt.Sprint(winner.fakeouts))

	return ret
}
