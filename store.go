package main

import (
	"sync"
	"time"
)

// Store holds active sessions with TTL-based cleanup.
type Store struct {
	sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
}

// NewStore returns an empty Store with the given TTL.
func NewStore(ttl time.Duration) *Store {
	return &Store{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
}

// Delete removes a session and kills the process (if running).
func (s *Store) Delete(uuid string) {
	s.Lock()
	if sess, ok := s.sessions[uuid]; ok {
		_ = sess.cmd.Process.Kill()
		delete(s.sessions, uuid)
	}
	s.Unlock()
}
