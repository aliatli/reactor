package main

import (
	"log"
	"net/http"

	"github.com/aliatli/reactor/internal/api"
	"github.com/aliatli/reactor/internal/db"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(database)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		log.Fatal(err)
	}
}
