package core

import "github.com/gorilla/websocket"

type Gamelike interface {
	Html() string
	NewHub() Hublike
	NewClient(hub Hublike, conn *websocket.Conn) Clientlike
}

type Game struct {
	html      string
	newHub    func() Hublike
	newClient func(hub Hublike, conn *websocket.Conn) Clientlike
}

func (g *Game) Html() string {
	return g.html
}

func (g *Game) NewHub() Hublike {
	return g.newHub()
}

func (g *Game) NewClient(hub Hublike, conn *websocket.Conn) Clientlike {
	return g.newClient(hub, conn)
}

func NewGame(html string, newHub func() Hublike, newClient func(hub Hublike, conn *websocket.Conn) Clientlike) *Game {
	return &Game{html: html, newHub: newHub, newClient: newClient}
}
