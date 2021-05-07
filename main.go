// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	hubs = map[string]*SpecializedHub{}
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	urlString := r.URL.String()
	hub, ok := hubs[urlString]
	if !ok {
		hub = newHub()
		go hub.run()
		hubs[urlString] = hub
	}
	serveWs(hub, w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path[:min(3, len(r.URL.Path))] == "/ws" {
		handleWs(w, r)
	} else {
		serveHome(w, r)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	specializedInit()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
