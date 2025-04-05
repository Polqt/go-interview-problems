package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockNode struct {
	delay      time.Duration
	shouldFail bool
	response   []byte
}

func TestFirstNodeSucceeds(t *testing.T) {
	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		mockNodes := map[string]mockNode{
			"node1": {delay: 10 * time.Millisecond, shouldFail: false, response: []byte("success-1")},
			"node2": {delay: 20 * time.Millisecond, shouldFail: false, response: []byte("success-2")},
		}

		mockNode, exists := mockNodes[node.Address]
		if !exists {
			return nil, errors.New("unknown node")
		}

		select {
		case <-time.After(mockNode.delay):
			if mockNode.shouldFail {
				return nil, ErrTaskFailed
			}
			return mockNode.response, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
	}

	result, err := ExecuteWithFailover(context.Background(), nodes, "task1", []byte("test-payload"))
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if string(result) != "success-1" {
		t.Errorf("Expected result 'success-1', got '%s'", string(result))
	}
}

func TestFirstNodeFailsSecondNodeSucceeds(t *testing.T) {
	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		mockNodes := map[string]mockNode{
			"node1": {delay: 10 * time.Millisecond, shouldFail: true, response: nil},
			"node2": {delay: 20 * time.Millisecond, shouldFail: false, response: []byte("success-2")},
		}

		mockNode, exists := mockNodes[node.Address]
		if !exists {
			return nil, errors.New("unknown node")
		}

		select {
		case <-time.After(mockNode.delay):
			if mockNode.shouldFail {
				return nil, ErrTaskFailed
			}
			return mockNode.response, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
	}

	result, err := ExecuteWithFailover(context.Background(), nodes, "task1", []byte("test-payload"))
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if string(result) != "success-2" {
		t.Errorf("Expected result 'success-2', got '%s'", string(result))
	}
}

func TestRespectsPriorityOrder(t *testing.T) {
	var executionOrder []string
	var mu sync.Mutex

	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		mu.Lock()
		executionOrder = append(executionOrder, node.Address)
		mu.Unlock()

		delay := time.Duration(10*node.Priority) * time.Millisecond
		time.Sleep(delay)
		return []byte("success-" + node.Address), nil
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node3", Priority: 3},
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
	}

	result, err := ExecuteWithFailover(context.Background(), nodes, "task1", []byte("test-payload"))
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if string(result) != "success-node1" {
		t.Errorf("Expected result from highest priority node 'success-node1', got '%s'", string(result))
	}

	if len(executionOrder) < 1 || executionOrder[0] != "node1" {
		t.Errorf("Expected node1 (highest priority) to be executed first, got execution order: %v", executionOrder)
	}
}

func TestAllNodesFail(t *testing.T) {
	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		return nil, ErrTaskFailed
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
	}

	_, err := ExecuteWithFailover(context.Background(), nodes, "task1", []byte("test-payload"))
	if err == nil {
		t.Fatal("Expected error when all nodes fail, got nil")
	}
}

func TestHandles500msTimeoutCorrectly(t *testing.T) {
	var nodeStartTimes sync.Map
	var mu sync.Mutex
	var executionOrder []string

	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		nodeStartTimes.Store(node.Address, time.Now())

		mu.Lock()
		executionOrder = append(executionOrder, node.Address)
		mu.Unlock()

		if node.Address == "node1" {
			select {
			case <-time.After(1 * time.Second):
				return []byte("success-but-slow"), nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if node.Address == "node2" {
			time.Sleep(100 * time.Millisecond)
			return []byte("success-fast"), nil
		}

		if node.Address == "node3" {
			time.Sleep(50 * time.Millisecond)
			return nil, ErrTaskFailed
		}

		return nil, errors.New("unknown node")
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
		{Address: "node3", Priority: 3},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteWithFailover(ctx, nodes, "task1", []byte("test-payload"))
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if string(result) != "success-fast" {
		t.Errorf("Expected fastest successful result 'success-fast', got '%s'", string(result))
	}

	node1Start, _ := nodeStartTimes.Load("node1")
	node2Start, _ := nodeStartTimes.Load("node2")
	if node2Start.(time.Time).Sub(node1Start.(time.Time)) > 550*time.Millisecond {
		t.Errorf("Expected node2 to start within 500ms after node1")
	}

	if elapsed > 650*time.Millisecond {
		t.Errorf("Expected to get result in under 650ms, took %v", elapsed)
	}
}

func TestRespectsContextCancellation(t *testing.T) {
	mockExecuteOnNode := func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
		select {
		case <-time.After(500 * time.Millisecond):
			return []byte("success"), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	executeOnNode = mockExecuteOnNode
	defer func() {
		executeOnNode = func(ctx context.Context, node Node, taskID string, payload []byte) ([]byte, error) {
			return nil, nil
		}
	}()

	nodes := []Node{
		{Address: "node1", Priority: 1},
		{Address: "node2", Priority: 2},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := ExecuteWithFailover(ctx, nodes, "task1", []byte("test-payload"))
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected context deadline exceeded error, got: %v", err)
	}
}

func TestHandlesEmptyNodesList(t *testing.T) {
	_, err := ExecuteWithFailover(context.Background(), []Node{}, "task1", []byte("test-payload"))
	if err != ErrNoNodesAvailable {
		t.Fatalf("Expected ErrNoNodesAvailable, got: %v", err)
	}
}
