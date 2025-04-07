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
	data  string
	err   error
	delay time.Duration
}

func (m *MockGetter) Get(ctx context.Context, address string) (string, error) {
	resp, exists := m.responses[address]
	if !exists {
		return "", errors.New("unknown address")
	}

	if resp.delay > 0 {
		select {
		case <-time.After(resp.delay):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	return resp.data, resp.err
}

func TestExecuteWithFailover(t *testing.T) {
	tests := []struct {
		name     string
		addresss []string
		response string
		err      error
		ttl      time.Duration
	}{
		{
			name:     "first node succeeds",
			addresss: []string{"success", "error"},
			response: "success",
			err:      nil,
			ttl:      time.Second,
		},
		{
			name:     "failover to second node",
			addresss: []string{"error", "success"},
			response: "success",
			err:      nil,
			ttl:      time.Second,
		},
		{
			name:     "timeout to second node",
			addresss: []string{"timeout", "success"},
			response: "success",
			err:      nil,
			ttl:      time.Second * 3,
		},
		{
			name:     "all nodes fail",
			addresss: []string{"error", "error"},
			response: "",
			err:      ErrRequestsFailed,
			ttl:      time.Second,
		},
		{
			name:     "context cancellation",
			addresss: []string{"timeout", "timeout"},
			response: "",
			err:      context.DeadlineExceeded,
			ttl:      time.Millisecond * 10, // Short timeout to trigger cancellation
		},
		{
			name:     "slow node eventually succeeds",
			addresss: []string{"slow-success"},
			response: "slow success",
			err:      nil,
			ttl:      time.Second,
		},
	}

	client := &MockGetter{
		responses: map[string]mockResponse{
			"success":      {"success", nil, 0},
			"error":        {"", ErrRequestsFailed, 0},
			"slow-success": {"slow success", nil, 600 * time.Millisecond},
			"timeout":      {"timeout", nil, 2 * time.Second},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.ttl)
			defer cancel()

			resp, err := RequestWithFailover(ctx, client, tt.addresss)
			if err != tt.err {
				t.Errorf("Unexpected error: %v", err)
			}
			if resp != tt.response {
				t.Errorf("Expected: %s, got: %s", tt.response, resp)
			}
		})
	}
}
