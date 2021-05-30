package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func (h *IdiotmouthHub) validWord(str string) int {
	for _, v := range h.usedWords {
		if v == str {
			return 2
		}
	}
	for _, v := range words {
		if v == str {
			return 0
		}
	}
	return 1
}

func (h *IdiotmouthHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.pass = false
	}
	h.wordsLeft = len(words)
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
	for client := range h.getAssertedClients() {
		if client.pass {
			count++
		}
	}
	return count*2 > len(h.clients)
}

func (h *IdiotmouthHub) gotIt(word string) int {
	h.resetPass()
	h.usedWords = append(h.usedWords, word)
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

func (h *IdiotmouthHub) getPrompt() string {
	ret := string(h.start) + "*" + string(h.end)
	bonus := h.getWorth()
	ret += ", worth " + fmt.Sprint(bonus) + " points. There are " + fmt.Sprint(h.letters[string(h.start)+string(h.end)]) + " possible words"
	return ret
}

func (h *IdiotmouthHub) getScores() string {
	keys := make([]*IdiotmouthClient, 0, len(h.clients))
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
