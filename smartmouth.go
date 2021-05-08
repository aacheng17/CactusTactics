package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	baseNumWords = 370000
)

var (
	words       []string
	letters     = [26]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
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
			h.sendData(client, byte('0'), []byte(client.name+": "+string(m.data)))
		}
		if len(m.data) >= 3 && m.data[0] == byte(h.start) && m.data[len(m.data)-1] == byte(h.end) && isWord(string(m.data)) {
			worth := h.getWorth()
			worth *= len(m.data) - 2
			m.client.score += worth
			for client := range h.clients {
				h.sendData(client, byte('0'), []byte(m.client.name+" earned "+fmt.Sprint((h.getWorth()))+"x"+fmt.Sprint(len(m.data)-2)+"="+fmt.Sprint(worth)+" points"))
				h.sendData(client, byte('0'), []byte("."))
				h.genNextLetters()
				h.sendData(m.client, byte('2'), []byte(h.getPrompt()))
			}
			h.sendData(m.client, byte('1'), []byte(h.getScores()))
		} else if string(m.data) == "pass" {
			m.client.pass = true
			if h.getMajorityPass() {
				h.resetPass()
				h.genNextLetters()
				for client := range h.clients {
					h.sendData(client, byte('0'), []byte("Letters passed, new letters generated"))
					h.sendData(client, byte('0'), []byte("."))
					h.sendData(client, byte('2'), []byte(h.getPrompt()))
				}
			}
		} else if string(m.data) == "restart" {
			h.reset()
			for client := range h.clients {
				h.sendData(client, byte('1'), []byte(h.getScores()))
				h.sendData(client, byte('2'), []byte(h.getPrompt()))
				h.sendData(client, byte('3'), []byte(""))
				h.sendData(m.client, byte('0'), []byte("Game restarted"))
				h.sendData(client, byte('0'), []byte("."))
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
		h.sendData(m.client, byte('1'), []byte(h.getScores()))
		h.sendData(m.client, byte('2'), []byte(h.getPrompt()))
	}
}

func buildWords() []string {
	// Open the file
	csvfile, err := os.Open("english.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	words := make([]string, baseNumWords)

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	i := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if i < baseNumWords {
			words[i] = record[0]
		} else {
			words = append(words, record[0])
		}
		i++
	}
	return words
}

func getRandomWord() string {
	return words[rand.Intn(len(words))]
}

func isWord(str string) bool {
	for _, v := range words {
		if v == str {
			return true
		}
	}
	return false
}

func specializedInit() {
	rand.Seed(time.Now().Unix())
	words = buildWords()
}
