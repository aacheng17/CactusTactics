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
	startFreq   = map[string]int{"a": 724, "b": 470, "c": 844, "d": 462, "e": 370, "f": 292, "g": 291, "h": 383, "i": 372, "j": 70, "k": 97, "l": 267, "m": 535, "n": 287, "o": 333, "p": 1036, "q": 49, "r": 410, "s": 1068, "t": 550, "u": 692, "v": 146, "w": 169, "x": 16, "y": 29, "z": 40}
	endFreq     = map[string]int{"a": 532, "l": 633, "i": 84, "m": 375, "k": 113, "f": 39, "n": 847, "c": 479, "e": 1874, "u": 22, "b": 20, "h": 203, "y": 1174, "s": 1102, "t": 632, "r": 656, "d": 678, "o": 97, "p": 94, "g": 273, "w": 30, "x": 34, "z": 6, "v": 2, "j": 1, "q": 0}
	currentWord = ""
)

func handleClientMessage(c *SpecializedClient, d []byte) {
	log.Println(string(d))
	c.hub.messages <- newMessage(c, byte(d[0]), d[1:])
}

func handleHubMessage(h *Hub, m *Message) {
	switch m.messageType {
	case byte('0'):
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(client.name+": "+string(m.data)))
		}
	case byte('1'):
		name := string(m.data)
		if m.client.name == "" {
			m.client.name = name
		}
		for client := range h.clients {
			h.sendData(client, byte('0'), []byte(name+" joined"))
		}
		h.sendData(m.client, byte('1'), []byte(fmt.Sprint(m.client.score)))
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

func specializedInit() {
	rand.Seed(time.Now().Unix())
	words = buildWords()
}
