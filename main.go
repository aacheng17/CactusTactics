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

	"example.com/hello/core"
	"example.com/hello/fakeout"
	"example.com/hello/idiotmouth"
	"example.com/hello/utility"
)

var (
	hubs = make(map[string]map[string]core.Hublike)
)

func getHtml(game string) string {
	switch game {
	case "idiotmouth":
		return "idiotmouth/idiotmouth.html"
	case "fakeout":
		return "fakeout/fakeout.html"
	}
	return ""
}

func getHubmaker(game string) func() core.Hublike {
	switch game {
	case "idiotmouth":
		return idiotmouth.NewIdiotmouthHub
	case "fakeout":
		return fakeout.NewFakeoutHub
	}
	return nil
}

func getClientmaker(game string) func(hub core.Hublike, conn *websocket.Conn) core.Clientlike {
	switch game {
	case "idiotmouth":
		return idiotmouth.NewIdiotmouthClient
	case "fakeout":
		return fakeout.NewFakeoutClient
	}
	return nil
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	game := utility.UrlIndexGetPath(r.URL.String(), 0)
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
	game := utility.UrlIndexGetPath(r.URL.String(), 0)
	html := getHtml(game)
	if html == "" {
		return
	}
	_, ok := hubs[game]
	if !ok {
		hubs[game] = make(map[string]core.Hublike)
	}

	hubId := utility.UrlIndexGetPath(r.URL.String(), 1)
	hub, ok := hubs[game][hubId]
	if !ok {
		hub = getHubmaker(game)()
		go hub.Run()
		hubs[game][hubId] = hub
	}

	core.ServeWs(hub, w, r, getClientmaker(game))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, found := utility.Find(r.Header["Connection"], "Upgrade")
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
	idiotmouth.IdiotmouthInit()
	fakeout.FakeoutInit()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
