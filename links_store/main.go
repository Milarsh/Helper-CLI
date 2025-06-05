package main

import (
	"links_store/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	store, err := internal.New()
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("links_store listening on %s", addr)

	if err := http.ListenAndServe(addr, internal.Routes(store)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
