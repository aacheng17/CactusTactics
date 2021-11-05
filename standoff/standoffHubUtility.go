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
		if client.Name == "" {
			continue
		}
		client.kills = nil
		client.active = true
		client.alive = true
		client.roundsAlive = 0
		client.decision = -1
	}
	h.round = 0
}

func (h *StandoffHub) nextRound() {
	h.phase = 0
	for client := range h.getAssertedClients() {
		client.decision = -1
	}
	h.round++
}

func (h *StandoffHub) getPrompt() []string {
	ret := []string{fmt.Sprint(h.round)}
	for client := range h.getAssertedClients() {
		if client.active && client.alive {
			ret = append(ret, fmt.Sprint(client.id))
			ret = append(ret, client.Name)
		}
	}
	return ret
}

func (h *StandoffHub) isAllDecided() bool {
	for client := range h.getAssertedClients() {
		if client.decision == -1 && client.active && client.alive {
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
		if client.decision == -1 {
			continue
		}
		if client.decision == -2 {
			ret = append(ret, fmt.Sprint(u.TagId("p", h.useMessageNum()), u.Tag("b")+client.Name+u.ENDTAG, " shot nobody", u.ENDTAG))
			continue
		}
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
				ret = append(ret, h.kill(client.id, client.Name))
				client.kills = append(client.kills, client.Name)
				client.alive = false
			}
		} else {
			found := false
			for client2 := range h.getAssertedClients() {
				if client2.id == client.decision {
					found = true
					if client2.decision == client2.id {
						break
					}
					ret = append(ret, h.kill(client.id, client2.Name))
					client.kills = append(client.kills, client2.Name)
					client2.alive = false
					break
				}
			}
			if !found {
				ret = append(ret, h.kill(client.id, ""))
			}
		}
	}
	return ret
}

func (h *StandoffHub) reflect(killer int, victims []string) string {
	for client := range h.getAssertedClients() {
		if client.id == killer {
			victimsString := ""
			for _, victim := range victims {
				victimsString += " " + victim
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
		if keys[i].active != keys[j].active {
			return keys[i].active
		}
		if keys[i].alive != keys[j].alive {
			return keys[i].alive
		}
		if keys[i].roundsAlive != keys[j].roundsAlive {
			return keys[i].roundsAlive > keys[j].roundsAlive
		}
		return len(keys[i].kills) > len(keys[j].kills)
	})
	players := []string{}
	for _, client := range keys {
		if client.Name == "" {
			continue
		}
		players = append(players, client.Name)
		players = append(players, fmt.Sprint(client.Avatar))
		players = append(players, fmt.Sprint(client.Color))
		status := "alive"
		if !client.active {
			status = "spectating"
		} else if !client.alive {
			status = "dead"
		}
		players = append(players, fmt.Sprint(status))
		players = append(players, fmt.Sprint(len(client.kills)))
		dotdotdotStatus := "none"
		if client.active && client.alive {
			if client.decision == -1 {
				dotdotdotStatus = "dotdotdot"
			} else {
				dotdotdotStatus = "ready"
			}
		}
		players = append(players, fmt.Sprint(dotdotdotStatus))
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
		if keys[i].active != keys[j].active {
			return keys[i].active
		}
		if keys[i].alive != keys[j].alive {
			return keys[i].alive
		}
		if keys[i].roundsAlive != keys[j].roundsAlive {
			return keys[i].roundsAlive > keys[j].roundsAlive
		}
		return len(keys[i].kills) > len(keys[j].kills)
	})

	if keys[0].alive {
		ret = append(ret, "1")
	} else {
		ret = append(ret, "0")
	}

	for _, key := range keys {
		if !key.active {
			continue
		}
		kills := ""
		for j, kill := range key.kills {
			kills += kill
			if j < len(key.kills)-1 {
				kills += ", "
			}
		}
		ret = append(ret, key.Name, fmt.Sprint(key.roundsAlive), kills)
	}

	return ret
}
