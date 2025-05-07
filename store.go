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

// Delete stops the session (if any) and returns true if it existed.
func (s *Store) Delete(uuid string) bool {
	s.Lock()
	defer s.Unlock()

	sess, ok := s.sessions[uuid]
	if ok {
		_ = sess.stdin.Close()
		_ = sess.cmd.Process.Kill()
		delete(s.sessions, uuid)
	}
	return ok
}
