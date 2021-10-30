package standoff

import (
	"fmt"
	"sort"

	u "example.com/hello/utility"
)

func (h *StandoffHub) useMessageNum() int {
	ret := h.messageNum
	h.messageNum++
	return ret
}

func (h *StandoffHub) reset() {
	for client := range h.getAssertedClients() {
		client.kills = nil
		client.active = true
	}
	h.round = 0
	h.nextRound()
}

func (h *StandoffHub) nextRound() {
	for client := range h.getAssertedClients() {
		h.phase = 0
		client.decision = -1
	}
	h.round++
}

func (h *StandoffHub) getPrompt() []string {
	ret := []string{fmt.Sprint(h.round)}
	for client := range h.getAssertedClients() {
		ret = append(ret, fmt.Sprint(client.id))
		ret = append(ret, client.Name)
	}
	return ret
}

func (h *StandoffHub) isAllDecided() bool {
	for client := range h.getAssertedClients() {
		if client.decision == -1 {
			return false
		}
	}
	return true
}

func (h *StandoffHub) numAlive() int {
	livingCount := 0
	for client := range h.getAssertedClients() {
		if client.active && client.alive {
			livingCount++
		}
	}
	return livingCount
}

func (h *StandoffHub) calcResult() []string {
	ret := []string{}
	for client := range h.getAssertedClients() {
		if client.decision == client.id {
			reflections := []string{}
			for client2 := range h.getAssertedClients() {
				if client2 == client {
					continue
				}
				if client2.decision == client.id {
					reflections = append(reflections, client2.Name)
					client.kills = append(client.kills, client2.Name)
					client2.alive = false
				}
			}
			if len(reflections) != 0 {
				ret = append(ret, h.reflect(client.id, reflections))
			} else {
				client.kills = append(client.kills, client.Name)
				client.alive = false
			}
		}
		found := false
		for client2 := range h.getAssertedClients() {
			if client2.id == client.decision {
				found = true
				ret = append(ret, h.kill(client.id, client2.Name))
				break
			}
		}
		if !found {
			ret = append(ret, h.kill(client.id, ""))
		}
	}
	return ret
}

func (h *StandoffHub) reflect(killer int, victims []string) string {
	for client := range h.getAssertedClients() {
		if client.id == killer {
			victimsString := ""
			for _, victim := range victims {
				victimsString += victim
			}
			return fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+client.Name+u.ENDTAG, " reflected "+victimsString, u.ENDTAG)
		}
	}
	return ""
}

func (h *StandoffHub) kill(killer int, victim string) string {
	for client := range h.getAssertedClients() {
		if client.id == killer {
			if victim == "" {
				return fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+client.Name+u.ENDTAG, " shot at someone who left the game", u.ENDTAG)
			}
			return fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+client.Name+u.ENDTAG, " shot "+victim, u.ENDTAG)
		}
	}
	return ""
}

func (h *StandoffHub) getPlayers(excepts ...*StandoffClient) []string {
	keys := make([]*StandoffClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		isExcept := false
		for _, e := range excepts {
			if k == e {
				isExcept = true
				break
			}
		}
		if !isExcept {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].alive || !keys[j].active
	})
	players := []string{}
	for _, client := range keys {
		if client.Name == "" {
			continue
		}
		players = append(players, client.Name)
		players = append(players, fmt.Sprint(client.Avatar))
		players = append(players, fmt.Sprint(client.Color))
		players = append(players, fmt.Sprint(client.active))
		players = append(players, fmt.Sprint(client.alive))
	}
	return players
}

func (h *StandoffHub) getWinners() []string {
	ret := []string{}
	keys := make([]*StandoffClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i].kills) > len(keys[j].kills)
	})

	for _, key := range keys {
		if !key.active {
			continue
		}
		kills := ""
		for _, kill := range key.kills {
			kills += " " + kill
		}
		alive := "DEAD"
		if key.alive {
			alive = "ALIVE"
		}
		ret = append(ret, fmt.Sprint(u.TagId("p", h.useMessageNum()), alive+": ", u.Tag("b")+key.Name+u.ENDTAG, " KILLS: "+kills, u.ENDTAG))
	}

	return ret
}
