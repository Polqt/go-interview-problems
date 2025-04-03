package main

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("key not found")

func get(ctx context.Context, address string, key string) (string, error) {
	// Already implemented
	return "", nil
}

// Call `get` function for each received address in parallel
// Return first response or an error if all requests fail
func Get(ctx context.Context, adresses []string, key string) (string, error) {
	return "", nil
}
