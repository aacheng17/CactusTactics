package fakeout

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var decks []Deck

type DeckRaw struct {
	Instructions string          `json:"instructions"`
	Questions    json.RawMessage `json:"questions"`
}

type Deck struct {
	Name         string
	Instructions string
	Questions    []Question
}

type Question struct {
	Category           string   `json:"category"`
	Question           string   `json:"question"`
	Answer             string   `json:"answer"`
	AlternateSpellings []string `json:"alternateSpellings"`
	Suggestions        []string `json:"suggestions"`
}

func buildQuestions() {
	jsonFile, err := os.Open("fakeout/fakeoutQuestions.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var decksRaw map[string]json.RawMessage
	json.Unmarshal(byteValue, &decksRaw)
	for key := range decksRaw {
		var deckRaw DeckRaw
		json.Unmarshal(decksRaw[key], &deckRaw)
		var deck Deck
		deck.Name = key
		deck.Instructions = deckRaw.Instructions
		json.Unmarshal(deckRaw.Questions, &deck.Questions)
		for i, x := range deck.Questions {
			deck.Questions[i].Answer = strings.ToLower(x.Answer)
			for j, y := range x.AlternateSpellings {
				deck.Questions[i].AlternateSpellings[j] = strings.ToLower(y)
			}
			for j, y := range x.Suggestions {
				deck.Questions[i].Suggestions[j] = strings.ToLower(y)
			}
		}
		decks = append(decks, deck)
		fmt.Println(fmt.Sprint("Fakeout deck [", key, "] with ", len(deck.Questions), " questions"))
	}
}

func getFakeoutDeckOptions() []string {
	deckOptions := make([]string, len(decks))
	i := 0
	for _, v := range decks {
		deckOptions[i] = v.Name
		i++
	}
	return deckOptions
}

func getFakeoutQuestion(deck int, n int) Question {
	return decks[deck].Questions[n]
}
