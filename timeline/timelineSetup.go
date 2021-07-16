package timeline

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	events map[string]Event
)

func buildEvents() {
	events = make(map[string]Event)
	file, err := os.Open("timeline/events.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		data := strings.Split(s, "\t")
		date, err := strconv.Atoi(data[0])
		if err == nil {
			event := data[1]
			events[event] = date
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
