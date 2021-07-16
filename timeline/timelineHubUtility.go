package timeline

import (
	"fmt"
	"math/rand"
	"sort"
)

func (h *TimelineHub) useMessageNum() int {
	ret := h.messageNum
	h.messageNum++
	return ret
}

func (h *TimelineHub) newPlayerInitiative() int {
	ret := h.playerNum
	h.playerNum++
	return ret
}

func (h *TimelineHub) reset() {
	for client := range h.getAssertedClients() {
	}
	h.messageNum = 0
	h.events = make(map[string]Event)
	for k, v := range events {
		h.events[k] = v
	}
}

func randMapKey(m map[string]Event) string {
	r := rand.Intn(len(m))
	for k := range m {
		if r == 0 {
			return k
		}
		r--
	}
	panic("unreachable")
}

func (h *TimelineHub) newPrompt() {
	k := randMapKey(h.events)
	h.event = h.events[k]
	delete(h.events, k)
}

func (h *TimelineHub) getPrompt() []string {
	return []string{h.event.title}
}

func (h *TimelineHub) getTurn() []string {
	ret := []string{}
	for k, v := range h.events {
		ret = append(ret, fmt.Sprint(v))
		ret = append(ret, k)
	}
	return ret
}

func (h *TimelineHub) getPlayers(excepts ...*TimelineClient) []string {
	keys := make([]*TimelineClient, 0, len(h.Clients))
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
		return keys[i].score > keys[j].score
	})
	players := []string{}
	for _, client := range keys {
		if client.Name == "" {
			continue
		}
		players = append(players, client.Name)
		players = append(players, fmt.Sprint(client.Avatar))
		players = append(players, fmt.Sprint((client.Color)))
		players = append(players, fmt.Sprint(client.score))
		players = append(players, client.highestWord)
		players = append(players, fmt.Sprint(client.highestScore))
	}
	return players
}

func (h *TimelineHub) getWinners() []string {
	ret := []string{}
	keys := make([]*TimelineClient, 0, len(h.Clients))
	for k := range h.getAssertedClients() {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].score > keys[j].score
	})
	winner := keys[0]
	ret = append(ret, winner.Name)
	ret = append(ret, fmt.Sprint(winner.score))

	return ret
}
