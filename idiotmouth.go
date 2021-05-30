package main

import (
	"math/rand"
	"time"
)

func idiotmouthInit() {
	rand.Seed(time.Now().Unix())
	buildDictionary()
}
