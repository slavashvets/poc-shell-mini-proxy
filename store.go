package main

import (
	"sync"
	"time"
)

// Store holds interactive sessions with TTL-based cleanup.
type Store struct {
	sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
}

// NewStore allocates an empty Store.
func NewStore(ttl time.Duration) *Store {
	return &Store{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
}

// Delete stops the shell and removes the session.
func (s *Store) Delete(uuid string) {
	s.Lock()
	if sess, ok := s.sessions[uuid]; ok {
		_ = sess.stdin.Close()
		_ = sess.cmd.Process.Kill()
		delete(s.sessions, uuid)
	}
	s.Unlock()
}
