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

	// logging middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	addr := ":8080"
	log.Printf("🚀 server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
