package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
*	TODO:
	  Later:
		- user pages
		- figure out an actual goal for the game
*
*/

func main() {

	// prepare router
	r := chi.NewRouter()

	// open hub and start game
	chatHub := newHub()
	go chatHub.run()
	go chatHub.runGameClock()

	// disconnect to DB on application exit
	defer DB.Disconnect(context.Background())

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
	fmt.Println("* Listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
