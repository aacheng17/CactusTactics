package core

import (
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

type Gamelike interface {
	Name() string
	ExecuteTemplate(w http.ResponseWriter)
	NewHub() Hublike
	NewClient(hub Hublike, conn *websocket.Conn) Clientlike
}

type Game struct {
	Child     Gamelike
	Html      string
	newHub    func() Hublike
	newClient func(hub Hublike, conn *websocket.Conn) Clientlike
}

func (g *Game) NewHub() Hublike {
	return g.newHub()
}

func (g *Game) NewClient(hub Hublike, conn *websocket.Conn) Clientlike {
	return g.newClient(hub, conn)
}

type TemplateData struct {
}

func (g *Game) ExecuteTemplate(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("static/index.html", "static/"+g.Child.Name()+"/stylesheets.html", "static/"+g.Child.Name()+"/ingame.html"))
	templateData := TemplateData{}
	tmpl.ExecuteTemplate(w, "index.html", templateData)
}

func NewGame(newHub func() Hublike, newClient func(hub Hublike, conn *websocket.Conn) Clientlike) *Game {
	return &Game{newHub: newHub, newClient: newClient}
}
