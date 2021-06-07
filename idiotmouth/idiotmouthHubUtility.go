package idiotmouth

import (
	"fmt"
	"math/rand"
	"sort"
)

func (h *IdiotmouthHub) validWord(str string) int {
	if _, ok := h.usedWords[str]; ok {
		return 2
	}
	if _, ok := dictionary[str]; ok {
		return 0
	}
	return 1
}

func (h *IdiotmouthHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.pass = false
		client.highestWord = ""
		client.highestScore = 0
	}
	h.usedWords = make(map[string]bool)
	h.whattedWords = make(map[string]bool)
	h.wordsLeft = len(dictionary)
	for k, v := range letters {
		h.letters[k] = v
	}
	h.genNextLetters()
}

func (h *IdiotmouthHub) resetPass() {
	for client := range h.getAssertedClients() {
		client.pass = false
	}
}

func (h *IdiotmouthHub) pass() int {
	h.resetPass()
	h.wordsLeft -= h.letters[string(h.start)+string(h.end)]
	h.letters[string(h.start)+string(h.end)] = 0
	return h.genNextLetters()
}

func (h *IdiotmouthHub) getMajorityPass() bool {
	count := 0
	clientsWithNames := 0
	for client := range h.getAssertedClients() {
		if client.Name != "" {
			clientsWithNames++
			if client.pass {
				count++
			}
		}
	}
	return count*2 > clientsWithNames
}

func (h *IdiotmouthHub) gotIt(word string) int {
	h.resetPass()
	h.usedWords[word] = true
	h.wordsLeft--
	h.letters[string(h.start)+string(h.end)]--
	return h.genNextLetters()
}

func (h *IdiotmouthHub) genNextLetters() int {
	if h.wordsLeft <= 0 {
		return 1
	}
	r := rand.Intn(h.wordsLeft)
	c := 0
	for lets, freq := range h.letters {
		c += freq
		if r < c {
			h.start = rune(lets[0])
			h.end = rune(lets[1])
			break
		}
	}
	return 0
}

func (h *IdiotmouthHub) getWorth() int {
	return int(50-50*(float32(letters[string(h.start)+string(h.end)]-minFreq)/float32(maxFreq-minFreq))) + 50
}

func (h *IdiotmouthHub) getPrompt() []string {
	return []string{string(h.start), string(h.end), fmt.Sprint(h.getWorth()), fmt.Sprint(h.letters[string(h.start)+string(h.end)])}
}

func (h *IdiotmouthHub) getPlayers() []string {
	keys := make([]*IdiotmouthClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
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
		players = append(players, fmt.Sprint(client.score))
		players = append(players, client.highestWord)
		players = append(players, fmt.Sprint(client.highestScore))
	}
	return players
}

func (h *IdiotmouthHub) getWinners() []string {
	ret := []string{}
	keys := make([]*IdiotmouthClient, 0, len(h.Clients))
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
		return keys[i].highestScore > keys[j].highestScore
	})
	winner = keys[0]
	ret = append(ret, winner.Name)
	ret = append(ret, winner.highestWord)
	ret = append(ret, fmt.Sprint(winner.highestScore))

	return ret
}
