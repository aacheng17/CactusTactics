package fakeout

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	questions Questions
)

type Questions struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	Category           string   `json:"category"`
	Question           string   `json:"question"`
	Answer             string   `json:"answer"`
	AlternateSpellings []string `json:"alternateSpellings"`
	Suggestions        []string `json:"suggestions"`
}

func fakeoutBuildQuestions() {
	jsonFile, err := os.Open("fakeout/fakeoutQuestions.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &questions)
	for i, x := range questions.Questions {
		questions.Questions[i].Answer = strings.ToLower(x.Answer)
		for j, y := range x.AlternateSpellings {
			questions.Questions[i].AlternateSpellings[j] = strings.ToLower(y)
		}
		for j, y := range x.Suggestions {
			questions.Questions[i].Suggestions[j] = strings.ToLower(y)
		}
	}
}

func (q Questions) size() int {
	return len(q.Questions)
}

func (q Questions) getQuestion(n int) Question {
	return q.Questions[n]
}
