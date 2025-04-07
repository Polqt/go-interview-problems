package main

import (
	"context"
	"errors"
)

type Client interface {
	Get(ctx context.Context, address string) (string, error)
}

var ErrRequestsFailed = errors.New("task execution failed")

// RequestWithFailover attempts to request a data from available addresses:
// 1. If error, immediately try the address without waiting
// 2. If an address doesn't respond within 500ms, try the next but keep the original request running
// 3. Return the first successful response, or all ErrRequestsFailed if all nodes fail
// 4. Properly handle context cancellation throughout the process
func RequestWithFailover(ctx context.Context, client Client, addresses []string) (string, error) {
	return "", nil
}
