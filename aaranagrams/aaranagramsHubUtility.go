package aaranagrams

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	u "example.com/hello/utility"
)

func (h *AaranagramsHub) useMessageNum() int {
	ret := h.messageNum
	h.messageNum++
	return ret
}

func (h *AaranagramsHub) handleWord(c *AaranagramsClient, indicesSelectedString string) string {
	indicesSelected := []int{}
	for _, c := range strings.Split(indicesSelectedString, ",") {
		if c == "" {
			continue
		}
		letterIndex, err := strconv.Atoi(c) // convert this to an integer
		if err != nil {
			return ""
		}
		indicesSelected = append(indicesSelected, letterIndex)
	}
	word := ""
	for _, letterIndex := range indicesSelected {
		word += string(h.letters[letterIndex])
	}
	word = strings.ToLower(word)
	isValidWord := h.isValidWord(word)
	if isValidWord != 0 {
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+c.Name+u.ENDTAG, " ", strings.ToUpper(word))})
	}
	switch isValidWord {
	case 0:
		// handle removal of letters
		for _, letterIndex := range indicesSelected {
			h.letters[letterIndex] = ' '
		}
		h.Broadcast(ToClientCode["LETTERS"], []string{string(h.letters)})

		worth := u.GetWordScore(word) * 15
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
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", mNum), u.Tag("b")+c.Name+u.ENDTAG, " earned ", worth, "x", bonus, "=", finalWorth, " points for ", strings.ToUpper(word), u.ENDTAG), strings.ToUpper(word)})
		h.Broadcast(ToClientCode["PLAYERS"], h.getPlayers())
		if c.score >= h.scoreToWin {
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p prebr", mNum), u.Tag("b")+c.Name+u.ENDTAG, " won the game!", u.ENDTAG)})
			h.endGame()
			break
		}
		if err == 1 {
			h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "All possible words have been used.", u.ENDTAG)})
			break
		}
	case 1:
		log.Println("Invalid word")
	case 2:
		h.Broadcast(ToClientCode["LOBBY_CHAT_MESSAGE"], []string{fmt.Sprint(u.TagId("p", h.useMessageNum()), "This word has already been used this game.", u.ENDTAG)})
	}
	return word
}

func (h *AaranagramsHub) isValidWord(word string) int {
	if _, ok := h.dictionary.usedWords[word]; ok {
		return 2
	}
	if _, ok := dictionary[word]; ok {
		return 0
	}
	return 1
}

func (h *AaranagramsHub) reset() {
	for i := range h.letters {
		h.letters[i] = ' '
	}
	for client := range h.getAssertedClients() {
		client.score = 0
		client.highestWord = ""
		client.highestScore = 0
	}
	h.messageNum = 0
	h.turn = 0
	h.dictionary.generate(h.minWordLength)
}

func (h *AaranagramsHub) gotAWord(word string) int {
	h.dictionary.usedWords[word] = true
	h.dictionary.wordsLeft--
	if h.dictionary.wordsLeft < 1 {
		return 1
	}
	return 0
}

func (h *AaranagramsHub) getChaosModeAsString() string {
	chaosModeString := "0"
	if h.chaosMode {
		chaosModeString = "1"
	}
	return chaosModeString
}

func (h *AaranagramsHub) getClientsSortedByJoinTime() []*AaranagramsClient {
	keys := make([]*AaranagramsClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].JoinTime < keys[j].JoinTime
	})
	return keys
}

func (h *AaranagramsHub) getClientOfCurrentTurn() *AaranagramsClient {
	return h.getClientsSortedByJoinTime()[h.turn]
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
		if keys[i].score != keys[i].score {
			return keys[i].score > keys[j].score
		}
		return keys[i].JoinTime < keys[j].JoinTime
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
		turn := ""
		if !h.chaosMode && h.phase == Phase["PLAY"] && h.getClientOfCurrentTurn() == client {
			turn = "turn"
		}
		players = append(players, fmt.Sprint(turn))
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
