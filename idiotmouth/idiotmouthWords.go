package idiotmouth

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var (
	dictionary map[string]string
)

func buildWords() {
	//dictionary source: https://boardgames.stackexchange.com/questions/38366/latest-collins-scrabble-words-list-in-text-file https://drive.google.com/file/d/1XIFdZukAcDRiDIOgR_rHpICrrgJbLBxV/view

	dictionary = make(map[string]string)
	file, err := os.Open("idiotmouth/dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" || s == "Collins Scrabble Words (2019). 279,496 words with definitions." {
			continue
		}
		data := strings.Split(s, "\t")
		word := strings.ToLower(data[0])
		definition := data[1]
		dictionary[word] = definition
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func buildDictionary() {
	buildWords()
}
