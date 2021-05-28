package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	words       []string
	letters     map[string]int
	minFreq     int
	maxFreq     int
	startFreq   = map[rune]int{'a': 724, 'b': 470, 'c': 844, 'd': 462, 'e': 370, 'f': 292, 'g': 291, 'h': 383, 'i': 372, 'j': 70, 'k': 97, 'l': 267, 'm': 535, 'n': 287, 'o': 333, 'p': 1036, 'q': 49, 'r': 410, 's': 1068, 't': 550, 'u': 692, 'v': 146, 'w': 169, 'x': 16, 'y': 29, 'z': 40}
	endFreq     = map[rune]int{'a': 532, 'l': 633, 'i': 84, 'm': 375, 'k': 113, 'f': 39, 'n': 847, 'c': 479, 'e': 1874, 'u': 22, 'b': 20, 'h': 203, 'y': 1174, 's': 1102, 't': 632, 'r': 656, 'd': 678, 'o': 97, 'p': 94, 'g': 273, 'w': 30, 'x': 34, 'z': 6, 'v': 2, 'j': 1, 'q': 0}
	currentWord = ""
)

func handleClientMessage(c *SpecializedClient, d []byte) {
	log.Println(string(d))
	c.hub.messages <- newMessage(c, byte(d[0]), d[1:])
}

func handleHubMessage(h *SpecializedHub, m *Message) {
	switch m.messageType {
	case byte('0'):
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(m.client.name+": "+string(m.data)))
		}
		word := strings.TrimSpace(strings.ToLower(string(m.data)))
		if len(word) >= 3 && word[0] == byte(h.start) && word[len(word)-1] == byte(h.end) {
			switch h.validWord(word) {
			case 0:
				worth := h.getWorth()
				bonus := len(word) - 2
				finalWorth := worth * bonus
				m.client.score += finalWorth
				err := h.gotIt(word)
				if err == 1 {
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte(m.client.name+" earned "+fmt.Sprint(worth)+"x"+fmt.Sprint(bonus)+"="+fmt.Sprint(finalWorth)+" points"))
					h.sendData(client, byte('0'), []byte("."))
					h.sendData(client, byte('2'), []byte(h.getPrompt()))
					h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
					h.sendData(client, byte('1'), []byte(h.getScores()))
				}
			case 2:
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte("This word has already been used this game."))
				}
			}
		} else if string(m.data) == "pass" {
			m.client.pass = true
			if h.getMajorityPass() {
				err := h.pass()
				if err == 1 {
					for client := range h.clients {
						h.sendData(client, byte('0'), []byte("You have passed or gotten all possible words. Type restart to restart the game."))
					}
					break
				}
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte("Letters passed, new letters generated"))
					h.sendData(client, byte('0'), []byte("."))
					h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
					h.sendData(client, byte('2'), []byte(h.getPrompt()))
				}
			}
		} else if string(m.data) == "restart" {
			h.reset()
			for client := range h.clients {
				h.sendData(client, byte('3'), []byte(""))
				h.sendData(client, byte('0'), []byte(m.client.name+" restarted the game"))
				h.sendData(client, byte('1'), []byte(h.getScores()))
				h.sendData(client, byte('2'), []byte(h.getPrompt()))
				h.sendData(client, byte('0'), []byte("."))
				h.sendData(client, byte('0'), []byte("New letters: "+h.getPrompt()))
			}
		}
	case byte('1'):
		name := string(m.data)
		if m.client.name == "" {
			m.client.name = name
		}
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(name+" joined"))
		}
		for client := range h.clients {
			h.sendData(client, byte('1'), []byte(h.getScores()))
		}
		h.sendData(m.client, byte('2'), []byte(h.getPrompt()))
		h.sendData(m.client, byte('0'), []byte("New letters: "+h.getPrompt()))
	}
}

func buildWords() {
	csvfile, err := os.Open("english.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	words = []string{}
	r := csv.NewReader(csvfile)

	i := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		word := strings.ToLower(strings.TrimSpace(record[0]))
		if len(word) < 3 {
			continue
		}
		words = append(words, word)
		i++
	}
}

func buildLetters() {
	letters = make(map[string]int)
	for _, word := range words {
		first := word[0]
		last := word[len(word)-1]
		l := string(first) + string(last)
		if _, ok := letters[l]; !ok {
			letters[l] = 1
		} else {
			letters[l]++
		}
	}
}

func buildFreqs() {
	values := make([]int, 0, len(letters))
	for _, v := range letters {
		values = append(values, v)
	}
	for i, e := range values {
		if i == 0 || e < minFreq {
			minFreq = e
		}
	}
	for i, e := range values {
		if i == 0 || e > maxFreq {
			maxFreq = e
		}
	}
}

func specializedInit() {
	rand.Seed(time.Now().Unix())
	buildWords()
	buildLetters()
	buildFreqs()
}
