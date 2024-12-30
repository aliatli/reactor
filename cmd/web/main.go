package main

import (
	"log"
	"net/http"

	"github.com/aliatli/reactor/internal/api"
)

func main() {
	server := api.NewServer()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server.router); err != nil {
		log.Fatal(err)
	}
}
