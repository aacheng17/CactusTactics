package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	jsonFile, err := os.Open("fakeoutQuestions.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &questions)
	/*
		for _, x := range questions.Questions {
			log.Println(x.Question)
		}
	*/
}

func (q Questions) size() int {
	return len(q.Questions)
}

func (q Questions) getQuestion(n int) Question {
	return q.Questions[n]
}
