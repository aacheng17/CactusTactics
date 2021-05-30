// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	hubs = make(map[string]map[string]Hublike)
)

func getHtml(game string) string {
	switch game {
	case "idiotmouth":
		return "idiotmouth.html"
	case "fakeout":
		return "fakeout.html"
	}
	return ""
}

func getHubmaker(game string) func() Hublike {
	switch game {
	case "idiotmouth":
		return newIdiotmouthHub
	case "fakeout":
		return newFakeoutHub
	}
	return nil
}

func getClientmaker(game string) func(hub Hublike, conn *websocket.Conn) Clientlike {
	switch game {
	case "idiotmouth":
		return newIdiotmouthClient
	case "fakeout":
		return newFakeoutClient
	}
	return nil
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	game := urlIndexGetPath(r.URL.String(), 0)
	if game == "" {
		//serve home
		return
	}
	html := getHtml(game)
	if html == "" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, html)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	game := urlIndexGetPath(r.URL.String(), 0)
	html := getHtml(game)
	if html == "" {
		return
	}
	_, ok := hubs[game]
	if !ok {
		hubs[game] = make(map[string]Hublike)
	}

	hubId := urlIndexGetPath(r.URL.String(), 1)
	hub, ok := hubs[game][hubId]
	if !ok {
		hub = getHubmaker(game)()
		go hub.run()
		hubs[game][hubId] = hub
	}

	serveWs(hub, w, r, getClientmaker(game))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, found := find(r.Header["Connection"], "Upgrade")
	if !found {
		servePage(w, r)
	} else {
		handleWs(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler)
	rand.Seed(time.Now().Unix())
	idiotmouthInit()
	fakeoutInit()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
