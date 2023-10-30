package aaranagrams

import (
	"fmt"
	"math/rand"
	"sort"

	u "example.com/hello/utility"
)

func (h *AaranagramsHub) useMessageNum() int {
	ret := h.messageNum
	h.messageNum++
	return ret
}

func (h *AaranagramsHub) handleWord(c *AaranagramsClient, word string) {
	switch h.isValidWord(word) {
	case 0:
		worth := h.getWorth()
		bonus := len(word)
		finalWorth := worth * bonus
		c.score += finalWorth
		if finalWorth > c.highestScore {
			c.highestWord = word
			c.highestScore = finalWorth
		}
		err := h.gotAWord(word)
		mNum := h.useMessageNum()
		h.dictionary.whattedWords[mNum] = word
		h.Broadcast(ToClientCode["MESSAGE_WITH_WHAT"], []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", word, u.ENDTAG), word})
		h.Broadcast(ToClientCode["PROMPT"], h.getPrompt())
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		if c.score >= h.scoreToWin {
			h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p prebr", mNum), u.Tag("b")+c.Name+u.ENDTAG, " won the game!", u.ENDTAG)})
			h.endGame()
			return
		}
		if err == 1 {
			h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used or passed.", u.ENDTAG)})
			break
		}
	case 2:
		h.Broadcast(ToClientCode["GAME_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "This word has already been used this game.", u.ENDTAG)})
	}
}

func (h *AaranagramsHub) isValidWord(word string) int {
	firstLetter := rune(word[0])
	lastLetter := rune(word[len(word)-1])
	if len(word) < h.minWordLength || firstLetter != h.start || lastLetter != h.end {
		return 1
	}
	if _, ok := h.dictionary.usedWords[word]; ok {
		return 2
	}
	if _, ok := dictionary[word]; ok {
		return 0
	}
	return 1
}

func (h *AaranagramsHub) reset() {
	for client := range h.getAssertedClients() {
		client.score = 0
		client.pass = false
		client.highestWord = ""
		client.highestScore = 0
	}
	h.messageNum = 0
	h.dictionary.generate(h.minWordLength)
	h.genNextLetters()
}

func (h *AaranagramsHub) resetPass() {
	for client := range h.getAssertedClients() {
		client.pass = false
	}
}

func (h *AaranagramsHub) pass() int {
	h.resetPass()
	h.dictionary.wordsLeft -= h.dictionary.letters[string(h.start)+string(h.end)]
	h.dictionary.letters[string(h.start)+string(h.end)] = 0
	return h.genNextLetters()
}

func (h *AaranagramsHub) getMajorityPass() bool {
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

func (h *AaranagramsHub) gotAWord(word string) int {
	h.resetPass()
	h.dictionary.usedWords[word] = true
	h.dictionary.wordsLeft--
	h.dictionary.letters[string(h.start)+string(h.end)]--
	return h.genNextLetters()
}

func (h *AaranagramsHub) genNextLetters() int {
	if h.dictionary.wordsLeft <= 0 {
		return 1
	}
	r := rand.Intn(h.dictionary.wordsLeft)
	c := 0
	for lets, freq := range h.dictionary.letters {
		c += freq
		if r < c {
			h.start = rune(lets[0])
			h.end = rune(lets[1])
			break
		}
	}
	return 0
}

func (h *AaranagramsHub) getWorth() int {
	return h.dictionary.getWorth(string(h.start) + string(h.end))
}

func (h *AaranagramsHub) getPrompt() []string {
	return []string{string(h.start), string(h.end), fmt.Sprint(h.getWorth()), fmt.Sprint(h.dictionary.letters[string(h.start)+string(h.end)])}
}

func (h *AaranagramsHub) getPlayers(excepts ...*AaranagramsClient) []string {
	keys := make([]*AaranagramsClient, 0, len(h.Clients))
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
		players = append(players, client.highestWord)
		players = append(players, fmt.Sprint(client.highestScore))
	}
	return players
}

func (h *AaranagramsHub) endGame() {
	h.Broadcast(ToClientCode["END_GAME"], h.getWinners())
	h.phase = Phase["PREGAME"]
}

func (h *AaranagramsHub) getWinners() []string {
	ret := []string{}
	keys := make([]*AaranagramsClient, 0, len(h.Clients))
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
	highestWorder := keys[0]
	ret = append(ret, highestWorder.Name)
	ret = append(ret, highestWorder.highestWord)
	ret = append(ret, fmt.Sprint(highestWorder.highestScore))

	return ret
}
