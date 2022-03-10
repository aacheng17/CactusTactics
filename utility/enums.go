package utility

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func GenerateEnums(name string, Phase map[string]byte, ToServerCode map[string]byte, ToClientCode map[string]byte) {
	generateEnumsFromFile("static/"+name+"/enum.js", name, Phase, ToServerCode, ToClientCode)
	generateEnumsFromFile("static/globalEnum.js", name, Phase, ToServerCode, ToClientCode)
}

func generateEnumsFromFile(fileName string, name string, Phase map[string]byte, ToServerCode map[string]byte, ToClientCode map[string]byte) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentEnum := ""
	for scanner.Scan() {
		rawText := scanner.Text()
		text := strings.TrimSpace(rawText)
		if len(rawText) < 1 || text == rawText {
			continue
		}
		if len(text) < 2 {
			continue
		}
		if text[len(text)-1:] == "{" {
			colonIndex := strings.Index(text, ":")
			if colonIndex == -1 {
				continue
			}
			currentEnum = text[:colonIndex]
			continue
		}
		firstQuoteIndex := strings.Index(text, "'")
		if firstQuoteIndex != -1 {
			colonIndex := strings.Index(text, ":")
			key := text[:colonIndex]
			value := text[firstQuoteIndex+1 : firstQuoteIndex+2]
			valueByte := byte([]rune(value)[0])
			if currentEnum == "Phase" {
				Phase[key] = valueByte
			} else if currentEnum == "ToServerCode" {
				ToServerCode[key] = valueByte
			} else if currentEnum == "ToClientCode" {
				ToClientCode[key] = valueByte
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
