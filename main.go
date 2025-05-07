package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const sessionTTL = 10 * time.Minute
	store := NewStore(sessionTTL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", store.Handle)

	addr := ":8080"
	log.Printf("🚀 server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
