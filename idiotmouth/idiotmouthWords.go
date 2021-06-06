package idiotmouth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"example.com/hello/utility"
)

var (
	dictionary map[string]string
	letters    map[string]int
	minFreq    int
	maxFreq    int
)

func buildWords() {
	jsonFile, err := os.Open("idiotmouth/dictionary_compact.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dictionaryRaw interface{}
	json.Unmarshal(byteValue, &dictionaryRaw)
	var dictionaryAsserted = dictionaryRaw.(map[string]interface{})
	dictionary = make(map[string]string)
	for k, v := range dictionaryAsserted {
		s := v.(string)
		s = utility.RemoveEscapes(s)
		dictionary[k] = s
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
