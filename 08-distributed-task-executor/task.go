package main

import (
	"context"
	"errors"
)

var (
	ErrTaskFailed       = errors.New("task execution failed")
	ErrNoNodesAvailable = errors.New("no nodes available")
)

// Represents a remote node that can execute tasks
type Node struct {
	Address  string
	Priority int // Lower is higher priority
}

// Executes a task on a single remote node
var executeOnNode func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
	// Already implemented
	return nil, nil
}

// ExecuteWithFailover attempts to execute a task on available nodes with the following requirements:
// 1. First try nodes in order of priority (lowest number first)
// 2. If a node fails, immediately try the next node without waiting
// 3. If a node doesn't respond within 500ms, try the next node but keep the original request running
// 4. Return the first successful result, or all errors if all nodes fail
// 5. Properly handle context cancellation throughout the process
func ExecuteWithFailover(ctx context.Context, nodes []Node, taskID string, payload []byte) ([]byte, error) {
	return nil, nil
}
