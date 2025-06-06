package main

import (
	"log"
	"net/http"
	"os"

	"articles/internal"
)

func main() {
	store, err := internal.New()
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	log.Printf("articles service listening on %s", addr)

	if err := http.ListenAndServe(addr, internal.Routes(store)); err != nil {
		log.Fatalf("server: %v", err)
	}
}
