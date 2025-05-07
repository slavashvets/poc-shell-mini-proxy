package main

import (
	"io"
	"strings"
)

// readCommand reads the entire request body and returns a trimmed string.
func readCommand(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
