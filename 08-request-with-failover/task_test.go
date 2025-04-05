package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

type MockGetter struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	data  []byte
	err   error
	delay time.Duration
}

func (m *MockGetter) Get(ctx context.Context, address string) ([]byte, error) {
	resp, exists := m.responses[address]
	if !exists {
		return nil, errors.New("unknown address")
	}

	if resp.delay > 0 {
		select {
		case <-time.After(resp.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resp.data, resp.err
}

func TestExecuteWithFailover(t *testing.T) {
	originalGetter := getter
	defer func() { getter = originalGetter }()

	mockGetter := &MockGetter{
		responses: map[string]mockResponse{
			"success":      {[]byte("success"), nil, 0},
			"error":        {nil, errors.New("failed"), 0},
			"slow-success": {[]byte("slow success"), nil, 600 * time.Millisecond},
			"timeout":      {[]byte("timeout"), nil, 2 * time.Second},
		},
	}
	getter = mockGetter

	t.Run("first node succeeds", func(t *testing.T) {
		result, err := RequestWithFailover(context.Background(), []string{"success", "error"})
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if string(result) != "success" {
			t.Errorf("Expected 'success', got: %s", string(result))
		}
	})

	t.Run("failover to second node", func(t *testing.T) {
		result, err := RequestWithFailover(context.Background(), []string{"error", "success"})
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if string(result) != "success" {
			t.Errorf("Expected 'success', got: %s", string(result))
		}
	})

	t.Run("timeout to second node", func(t *testing.T) {
		result, err := RequestWithFailover(context.Background(), []string{"timeout", "success"})
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if string(result) != "success" {
			t.Errorf("Expected 'success', got: %s", string(result))
		}
	})

	t.Run("all nodes fail", func(t *testing.T) {
		_, err := RequestWithFailover(context.Background(), []string{"error", "error"})

		if !errors.Is(err, ErrTaskFailed) {
			t.Fatalf("Expected ErrTaskFailed, got: %v", err)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		_, err := RequestWithFailover(ctx, []string{"timeout", "timeout"})

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled, got: %v", err)
		}
	})

	t.Run("slow node eventually succeeds", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		result, err := RequestWithFailover(ctx, []string{"slow-success"})
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if string(result) != "slow success" {
			t.Errorf("Expected 'slow success', got: %s", string(result))
		}
	})
}
