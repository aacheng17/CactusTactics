package timeline

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var (
	dictionary map[string]string
	letters    map[string]int
	minFreq    int
	maxFreq    int
)

func buildWords() {
	//dictionary source: https://boardgames.stackexchange.com/questions/38366/latest-collins-scrabble-words-list-in-text-file https://drive.google.com/file/d/1XIFdZukAcDRiDIOgR_rHpICrrgJbLBxV/view

	dictionary = make(map[string]string)
	file, err := os.Open("timeline/dictionary.txt")
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

func buildLetters() {
	letters = make(map[string]int)
	for word := range dictionary {
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
