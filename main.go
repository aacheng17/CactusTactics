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

	"example.com/hello/core"
	"example.com/hello/fakeout"
	"example.com/hello/idiotmouth"
	"example.com/hello/timeline"
	"example.com/hello/utility"
)

var (
	games = make(map[string]core.Gamelike)
	hubs  = make(map[string]map[string]core.Hublike)
)

func Init() {
	games["idiotmouth"] = idiotmouth.Init()
	games["fakeout"] = fakeout.Init()
	games["timeline"] = timeline.Init()
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	game := utility.UrlIndexGetPath(r.URL.String(), 0)
	log.Println(game)
	if game == "" {
		http.ServeFile(w, r, "static/home.html")
		return
	}
	_, ok := games[game]
	if !ok {
		http.Error(w, "404 Page Not Found. This is the custom 404 page.", http.StatusNotFound)
		return
	}
	games[game].ExecuteTemplate(w)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	game := utility.UrlIndexGetPath(r.URL.String(), 0)
	_, ok := games[game]
	if !ok {
		return
	}
	_, ok = hubs[game]
	if !ok {
		hubs[game] = make(map[string]core.Hublike)
	}

	hubId := utility.UrlIndexGetPath(r.URL.String(), 1)
	log.Println(hubId)
	hub, ok := hubs[game][hubId]
	if !ok {
		hub = games[game].NewHub()
		go hub.Run()
		hubs[game][hubId] = hub
	}

	core.ServeWs(hub, w, r, games[game].NewClient)
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
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)
	rand.Seed(time.Now().Unix())
	Init()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
