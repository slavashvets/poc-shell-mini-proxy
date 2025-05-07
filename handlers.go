package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Handle is the entry-point wired in main.go.
func (s *Store) Handle(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/")
	if uuid == "" {
		http.Error(w, "missing session ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.createHandler(w, r, uuid)
	case http.MethodGet:
		s.streamHandler(w, r, uuid)
	case http.MethodDelete:
		s.deleteHandler(w, r, uuid)
	default:
		w.Header().Set("Allow", "GET, POST, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

/* POST /{uuid} */
func (s *Store) createHandler(w http.ResponseWriter, r *http.Request, uuid string) {
	cmd, err := readCommand(r.Body)
	if err != nil || cmd == "" {
		http.Error(w, "invalid or empty command", http.StatusBadRequest)
		return
	}

	s.Lock()
	if _, exists := s.sessions[uuid]; exists {
		s.Unlock()
		http.Error(w, "session already exists", http.StatusConflict)
		return
	}
	sess, err := newSession(cmd)
	if err != nil {
		s.Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.sessions[uuid] = sess
	s.Unlock()

	time.AfterFunc(s.ttl, func() { s.Delete(uuid) })
	w.WriteHeader(http.StatusAccepted)
}

/* GET /{uuid} */
func (s *Store) streamHandler(w http.ResponseWriter, r *http.Request, uuid string) {
	s.RLock()
	sess, ok := s.sessions[uuid]
	s.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	for {
		select {
		case line, open := <-sess.out:
			if !open {
				return
			}
			fmt.Fprintf(w, "data:%s\n\n", line)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

/* DELETE /{uuid} */
func (s *Store) deleteHandler(w http.ResponseWriter, r *http.Request, uuid string) {
	s.Delete(uuid)
	w.WriteHeader(http.StatusNoContent)
}
