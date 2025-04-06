package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Response represents a mock response with optional error and delay
type Response struct {
	Value string
	Error error
	Delay time.Duration
}

// MockGetter implements the Getter interface for testing
type MockGetter struct {
	Responses map[string]map[string]Response
}

func NewMockGetter(responses map[string]map[string]Response) *MockGetter {
	return &MockGetter{
		Responses: responses,
	}
}

func (m *MockGetter) Get(ctx context.Context, address, key string) (string, error) {
	if responses, exists := m.Responses[address]; exists {
		if resp, keyExists := responses[key]; keyExists {
			// Simulate delay if set
			if resp.Delay > 0 {
				select {
				case <-time.After(resp.Delay):
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}

			// Return the error if set
			if resp.Error != nil {
				return "", resp.Error
			}

			// Return the response value
			return resp.Value, nil
		}
	}

	return "", errors.New("key not found")
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		responses map[string]map[string]Response
		addresses []string
		key       string
		ttl       time.Duration
		wantValue string
		wantErr   bool
	}{
		{
			name: "first address fails second succeeds",
			responses: map[string]map[string]Response{
				"addr1": {
					"key1": {Error: errors.New("connection error")},
				},
				"addr2": {
					"key1": {Value: "value2"},
				},
			},
			addresses: []string{"addr1", "addr2"},
			key:       "key1",
			wantValue: "value2",
			ttl:       1 * time.Millisecond,
			wantErr:   false,
		},
		{
			name: "all addresses fail",
			responses: map[string]map[string]Response{
				"addr1": {
					"key1": {Error: errors.New("error 1")},
				},
				"addr2": {
					"key1": {Error: errors.New("error 2")},
				},
			},
			addresses: []string{"addr1", "addr2"},
			key:       "key1",
			ttl:       1 * time.Millisecond,
			wantValue: "",
			wantErr:   true,
		},
		{
			name: "context cancellation",
			responses: map[string]map[string]Response{
				"addr1": {
					"key1": {Value: "value1", Delay: 200 * time.Millisecond},
				},
			},
			addresses: []string{"addr1"},
			key:       "key1",
			ttl:       50 * time.Millisecond,
			wantValue: "",
			wantErr:   true,
		},
		{
			name: "fast address wins over slow",
			responses: map[string]map[string]Response{
				"addr1": {
					"key1": {Value: "value1", Delay: 200 * time.Millisecond},
				},
				"addr2": {
					"key1": {Value: "value2", Delay: 50 * time.Millisecond},
				},
			},
			addresses: []string{"addr1", "addr2"},
			key:       "key1",
			ttl:       100 * time.Millisecond,
			wantValue: "value2",
			wantErr:   false,
		},
		{
			name:      "empty address list",
			responses: map[string]map[string]Response{},
			addresses: []string{},
			key:       "key1",
			ttl:       50 * time.Millisecond,
			wantValue: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockGetter(tt.responses)

			var getter Getter
			if tt.name == "nil getter" {
				getter = nil
			} else {
				getter = mock
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.ttl)
			got, err := Get(ctx, getter, tt.addresses, tt.key)
			cancel()

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantValue {
				t.Errorf("Get() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}
