package core

import (
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

type Gamelike interface {
	Name() string
	ExecuteTemplate(w http.ResponseWriter)
	NewHub(game string, id string, deleteHubCallback func(*Hub)) Hublike
	NewClient(hub Hublike, conn *websocket.Conn) Clientlike
}

type Game struct {
	Child     Gamelike
	Html      string
	newHub    func(game string, id string, deleteHubCallback func(*Hub)) Hublike
	newClient func(hub Hublike, conn *websocket.Conn) Clientlike
}

func (g *Game) NewHub(game string, id string, deleteHubCallback func(*Hub)) Hublike {
	return g.newHub(game, id, deleteHubCallback)
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

func NewGame(newHub func(game string, id string, deleteHubCallback func(*Hub)) Hublike, newClient func(hub Hublike, conn *websocket.Conn) Clientlike) *Game {
	return &Game{newHub: newHub, newClient: newClient}
}
