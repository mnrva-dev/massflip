package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

/*
	0.0.4
		- added rate limiting for http and ws
		- added captcha for account creation
		- added WebSocket authentication
		- fixed BetInput bugs (NaN, decimals)
		- seperated development and production environment
		- frontend bug fixes
*/

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("* No .env file found")
	}
	DB = openDB()
}

func main() {

	// prepare router
	r := chi.NewRouter()

	// open hub and start game
	chatHub := newHub()
	go chatHub.run()
	go chatHub.runGameClock()

	// disconnect to DB on application exit
	defer DB.Disconnect(context.Background())

	// rate limiting middleware
	r.Use(limitMiddleware)
	// handlers
	r.Handle("/*", http.FileServer(http.Dir("./frontend/dist")))
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatHub, w, r)
	})
	r.Post("/api/createaccount", createAccount)
	r.Post("/api/login", login)
	r.Post("/api/login/bysession", loginBySession)
	r.Post("/api/logout", logout)
	r.Put("/api/chatcolor", chatColor)

	// run server
	log.Println("* Listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
