package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// Session is a long-lived interactive shell.
type Session struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
	out   chan string
	done  chan struct{}
}

// newSession starts an interactive /bin/sh and returns a Session.
func newSession() (*Session, error) {
	cmd := exec.Command("/bin/sh") // interactive shell

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start shell: %w", err)
	}

	sess := &Session{
		cmd:   cmd,
		stdin: stdin,
		out:   make(chan string, 100),
		done:  make(chan struct{}),
	}

	go func() {
		defer close(sess.out)
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			sess.out <- scanner.Text()
		}
		cmd.Wait()
		close(sess.done)
	}()

	return sess, nil
}
