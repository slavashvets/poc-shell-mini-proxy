package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Handle dispatches PUT / POST / GET / DELETE for /{uuid}.
func (s *Store) Handle(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/")
	if uuid == "" {
		http.Error(w, "missing session ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		s.createHandler(w, r, uuid) // create session
	case http.MethodPost:
		s.execHandler(w, r, uuid) // run command
	case http.MethodGet:
		s.streamHandler(w, r, uuid) // SSE
	case http.MethodDelete:
		s.deleteHandler(w, r, uuid) // kill session
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

/* PUT /{uuid} — create session */
func (s *Store) createHandler(w http.ResponseWriter, _ *http.Request, uuid string) {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.sessions[uuid]; exists {
		http.Error(w, "session already exists", http.StatusConflict)
		return
	}
	sess, err := newSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.sessions[uuid] = sess
	time.AfterFunc(s.ttl, func() { s.Delete(uuid) })
	w.WriteHeader(http.StatusCreated)
}

/* POST /{uuid} — execute command within an existing shell */
func (s *Store) execHandler(w http.ResponseWriter, r *http.Request, uuid string) {
	cmdStr, err := readCommand(r.Body)
	if err != nil || cmdStr == "" {
		http.Error(w, "empty command", http.StatusBadRequest)
		return
	}

	s.RLock()
	sess, ok := s.sessions[uuid]
	s.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}

	// send command + newline to shell stdin
	if _, err := fmt.Fprintln(sess.stdin, cmdStr); err != nil {
		http.Error(w, "write to session failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

/* GET /{uuid} — SSE stream */
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

/* DELETE /{uuid} — terminate session */
func (s *Store) deleteHandler(w http.ResponseWriter, r *http.Request, uuid string) {
	s.Delete(uuid)
	w.WriteHeader(http.StatusNoContent)
}
