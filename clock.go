package main

import (
	"math/rand"
	"time"
)

const (
	// time until next flip
	GAME_TIME = 180
)

func (h *Hub) runGameClock() {
	var msg string
	for {
		for i := GAME_TIME; i >= 0; i-- {
			h.tick(i)
			time.Sleep(time.Second)
		}
		HorT := rand.Int() % 2
		if HorT == 1 {
			msg = "{\"type\":\"flip\",\"value\":\"heads\"}"
		} else {
			msg = "{\"type\":\"flip\",\"value\":\"tails\"}"
		}
		h.broadcast <- []byte(msg)
		time.Sleep(5 * time.Second)
		h.broadcastPoolUpdate()
	}
}
