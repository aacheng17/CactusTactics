package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

var (
	words     []string
	letters   map[string]int
	minFreq   int
	maxFreq   int
	startFreq = map[rune]int{'a': 724, 'b': 470, 'c': 844, 'd': 462, 'e': 370, 'f': 292, 'g': 291, 'h': 383, 'i': 372, 'j': 70, 'k': 97, 'l': 267, 'm': 535, 'n': 287, 'o': 333, 'p': 1036, 'q': 49, 'r': 410, 's': 1068, 't': 550, 'u': 692, 'v': 146, 'w': 169, 'x': 16, 'y': 29, 'z': 40}
	endFreq   = map[rune]int{'a': 532, 'l': 633, 'i': 84, 'm': 375, 'k': 113, 'f': 39, 'n': 847, 'c': 479, 'e': 1874, 'u': 22, 'b': 20, 'h': 203, 'y': 1174, 's': 1102, 't': 632, 'r': 656, 'd': 678, 'o': 97, 'p': 94, 'g': 273, 'w': 30, 'x': 34, 'z': 6, 'v': 2, 'j': 1, 'q': 0}
)

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

func buildDictionary() {
	buildWords()
	buildLetters()
	buildFreqs()
}
