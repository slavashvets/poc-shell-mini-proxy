package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Session represents a running shell process exposed via SSE.
type Session struct {
	cmd  *exec.Cmd
	out  chan string
	done chan struct{}
}

// readCommand reads and trims the request body.
func readCommand(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// newSession starts /bin/sh -c <command> and returns a Session.
func newSession(command string) (*Session, error) {
	cmd := exec.Command("/bin/sh", "-c", command) // ← контекст больше не нужен

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start command: %w", err)
	}

	sess := &Session{
		cmd:  cmd,
		out:  make(chan string, 100),
		done: make(chan struct{}),
	}

	go func() {
		defer close(sess.out)
		reader := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			sess.out <- scanner.Text()
		}
		cmd.Wait() // wait for process completion
		close(sess.done)
	}()

	return sess, nil
}
