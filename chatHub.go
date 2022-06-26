package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Chat Hub code adapted from https://github.com/gorilla/websocket/tree/master/examples/chat

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Betting Hub Data
	headsPool    int
	headsBetters map[string]int
	tailsPool    int
	tailsBetters map[string]int
	allUsers     map[string]string
	largestBet   int
}

type WsMessage struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
}

type WsBetMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Bet      string `json:"bet"`
	Amount   int    `json:"amount"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		headsBetters: make(map[string]int),
		tailsBetters: make(map[string]int),
		allUsers:     make(map[string]string),
	}
}

func (h *Hub) broadcastPoolUpdate() {
	var msg string
	heads := strconv.Itoa(len(h.headsBetters))
	tails := strconv.Itoa(len(h.tailsBetters))
	headspool := strconv.Itoa(h.headsPool)
	tailspool := strconv.Itoa(h.tailsPool)
	msg = "{\"type\":\"pool\",\"heads\":" + heads + ",\"tails\":" + tails + ",\"headspool\":" + headspool + ",\"tailspool\":" + tailspool + "}"

	h.broadcast <- []byte(msg)
}

func (h *Hub) tick(s int) {
	var msg string
	heads := strconv.Itoa(len(h.headsBetters))
	tails := strconv.Itoa(len(h.tailsBetters))
	headspool := strconv.Itoa(h.headsPool)
	tailspool := strconv.Itoa(h.tailsPool)
	msg = "{\"type\":\"tick\",\"clock\":" + strconv.Itoa(s) + ",\"heads\":" + heads + ",\"tails\":" + tails + ",\"headspool\":" + headspool + ",\"tailspool\":" + tailspool + "}"

	h.broadcast <- []byte(msg)
}

func (h *Hub) processMessage(message []byte) error {
	var msg WsMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		// need error handling
		fmt.Println("* Error unmarshalling WebSocket message")
		log.Println(string(message))
	}
	if msg.Type == "bet" {
		var betMsg WsBetMessage
		err := json.Unmarshal(message, &betMsg)
		if err != nil {
			fmt.Println("* Error unmarshalling WebSocket bet message")
		}
		if h.allUsers[betMsg.Username] != "" {
			log.Println("* Disallowed bet from user " + betMsg.Username + ": You can't bet twice")
			return fmt.Errorf("cant bet twice")
		}
		if betMsg.Amount > DBGetUserByUsername(betMsg.Username).Points || betMsg.Amount <= 0 {
			log.Println("* Disallowed bet from user " + betMsg.Username + ": Cannot bet more gp than you have or bet 0")
			return fmt.Errorf("cant bet more gp than you have")
		}
		err = DBSubtractPoints(betMsg.Username, betMsg.Amount)
		if err != nil {
			log.Println("* Error subtracting points from user " + betMsg.Username + " in DB")
		}
		if betMsg.Bet == "heads" {
			h.allUsers[betMsg.Username] = "heads"
			h.headsBetters[betMsg.Username] = betMsg.Amount
			h.headsPool += betMsg.Amount
			fmt.Println("* " + betMsg.Username + " has bet " + strconv.Itoa(h.headsBetters[betMsg.Username]) + " on heads")
		} else if betMsg.Bet == "tails" {
			h.allUsers[betMsg.Username] = "tails"
			h.tailsBetters[betMsg.Username] = betMsg.Amount
			h.tailsPool += betMsg.Amount
			fmt.Println("* " + betMsg.Username + " has bet " + strconv.Itoa(h.tailsBetters[betMsg.Username]) + " on tails")
		}
		h.broadcastPoolUpdate()
	} else if msg.Type == "flip" {

		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}

		// this 50ms delay prevents the WS messages from colliding when they reach
		// the client and causing a json parsing error
		time.Sleep(50 * time.Millisecond)

		var v struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err := json.Unmarshal(message, &v)
		if err != nil {
			fmt.Println("* Error unmarshalling WebSocket bet message")
		}

		if h.headsPool+h.tailsPool == 0 {
			//fmt.Println("* No betters in pool, skipping payout")
			return fmt.Errorf("no betters in pool")
		}
		for c := range h.clients {
			if h.headsBetters[c.username] > 0 && v.Value == "heads" {
				winAmt := int(float64(h.headsBetters[c.username])/float64(h.headsPool)*float64(h.tailsPool)) + h.headsBetters[c.username]
				m := "{\"type\":\"win\",\"value\":" + strconv.Itoa(winAmt) + "}"
				c.send <- []byte(m)
				DBAddPoints(c.username, winAmt)
				fmt.Println("* " + strconv.Itoa(winAmt) + " paid out to " + c.username)
			} else if h.tailsBetters[c.username] > 0 && v.Value == "tails" {
				winAmt := int(float64(h.tailsBetters[c.username])/float64(h.tailsPool)*float64(h.headsPool)) + h.tailsBetters[c.username]
				m := "{\"type\":\"win\",\"value\":" + strconv.Itoa(winAmt) + "}"
				c.send <- []byte(m)
				DBAddPoints(c.username, winAmt)
				fmt.Println("* " + strconv.Itoa(winAmt) + " paid out to " + c.username)
			}
		}
		h.headsBetters = make(map[string]int)
		h.tailsBetters = make(map[string]int)
		h.allUsers = make(map[string]string)
		h.headsPool = 0
		h.tailsPool = 0
		h.largestBet = 0
		// reset pool for clients
	} else {
		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
	return nil
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			go h.processMessage(message)
		}
	}
}
