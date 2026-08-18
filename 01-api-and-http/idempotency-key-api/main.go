package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewServer()

	log.Printf("server running on http://localhost:8000")

	if err := http.ListenAndServe(":8000", server.routes()); err != nil {
		log.Fatal(err)
	}
}
