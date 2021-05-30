package main

import (
	"math/rand"
	"time"
)

// MESSAGE TYPES (SERVER TO CLIENT)
// 0: regular chat messages
// 1: scores
// 2: prompt
// 3: restart (data is inconsequential, probably empty string)

func specializedInit() {
	rand.Seed(time.Now().Unix())
	buildDictionary()
}
