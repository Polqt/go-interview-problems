package main

import (
	"context"
	"errors"
	"sort"
	"time"
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
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Priority < nodes[j].Priority })

	nodeCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	ch := make(chan []byte, 1)
	errCh := make(chan error, len(nodes))

	var cur int
	go startNext(ctx, nodes[cur], taskID, payload, ch, errCh)

	var errJoin error
	for {
		select {
		case resp := <-ch:
			return resp, nil
		case err := <-errCh:
			errJoin = errors.Join(errJoin, err)
			cur++
			if cur < len(nodes) {
				nodeCtx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
				defer cancel()
				go startNext(ctx, nodes[cur], taskID, payload, ch, errCh)
				continue
			}
			return nil, errJoin
		case <-nodeCtx.Done():
			cur++
			if cur < len(nodes) {
				nodeCtx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
				defer cancel()
				go startNext(ctx, nodes[cur], taskID, payload, ch, errCh)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func startNext(ctx context.Context, node Node, taskID string, payload []byte, ch chan<- []byte, errCh chan<- error) {
	resp, err := executeOnNode(ctx, node, taskID, payload)
	if err != nil {
		errCh <- err
		return
	}
	select {
	case ch <- resp:
	default:
		return
	}
}
