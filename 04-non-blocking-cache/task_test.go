package main

import (
	"errors"
	"testing"
	"time"
)

var (
	ErrNoResponse = errors.New("no response")
	ErrExpected   = errors.New("expected error")
)

type response struct {
	body  string
	err   error
	delay time.Duration
}

type mockClient struct {
	responses map[string]chan response
}

func newMockClient(responses map[string][]response) *mockClient {
	client := &mockClient{responses: map[string]chan response{}}
	for address, responses := range responses {
		client.responses[address] = make(chan response, len(responses))
		for _, resp := range responses {
			client.responses[address] <- resp
		}
	}
	return client
}

func (c *mockClient) Get(address string) (string, error) {
	select {
	case resp := <-c.responses[address]:
		<-time.After(resp.delay)
		return resp.body, resp.err
	default:
		return "", ErrNoResponse
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		responses map[string][]response
		requests  []string
		results   []struct {
			body string
			err  error
		}
	}{
		{
			name: "Basic caching",
			responses: map[string][]response{
				"example.com": {
					{body: "response1", err: nil, delay: 100 * time.Millisecond},
				},
			},
			requests: []string{"example.com", "example.com"},
			results: []struct {
				body string
				err  error
			}{
				{body: "response1", err: nil},
				{body: "response1", err: nil},
			},
		},
		{
			name: "Multiple different URLs",
			responses: map[string][]response{
				"example.com": {
					{body: "response1", err: nil, delay: 50 * time.Millisecond},
				},
				"example.org": {
					{body: "response2", err: nil, delay: 50 * time.Millisecond},
				},
			},
			requests: []string{"example.com", "example.org", "example.com", "example.org"},
			results: []struct {
				body string
				err  error
			}{
				{body: "response1", err: nil},
				{body: "response2", err: nil},
				{body: "response1", err: nil},
				{body: "response2", err: nil},
			},
		},
		{
			name: "Error handling",
			responses: map[string][]response{
				"example.com": {
					{body: "", err: ErrExpected, delay: 50 * time.Millisecond},
				},
			},
			requests: []string{"example.com", "example.com"},
			results: []struct {
				body string
				err  error
			}{
				{body: "", err: ErrExpected},
				{body: "", err: ErrExpected},
			},
		},
		{
			name:      "URL not found",
			responses: map[string][]response{},
			requests:  []string{"nonexistent.com"},
			results: []struct {
				body string
				err  error
			}{
				{body: "", err: ErrNoResponse},
			},
		},
		{
			name: "Mixed success and errors",
			responses: map[string][]response{
				"success.com": {
					{body: "success", err: nil, delay: 50 * time.Millisecond},
				},
				"error.com": {
					{body: "", err: ErrExpected, delay: 50 * time.Millisecond},
				},
			},
			requests: []string{"success.com", "error.com", "success.com", "error.com"},
			results: []struct {
				body string
				err  error
			}{
				{body: "success", err: nil},
				{body: "", err: ErrExpected},
				{body: "success", err: nil},
				{body: "", err: ErrExpected},
			},
		},
		{
			name: "Changing responses",
			responses: map[string][]response{
				"example.com": {
					{body: "first response", err: nil, delay: 50 * time.Millisecond},
					{body: "second response", err: nil, delay: 50 * time.Millisecond},
				},
			},
			requests: []string{"example.com", "example.com", "example.com"},
			results: []struct {
				body string
				err  error
			}{
				{body: "first response", err: nil},
				{body: "first response", err: nil},
				{body: "first response", err: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockClient(tt.responses)
			cache := NewCache(client)
			for i, req := range tt.requests {
				resp, err := cache.Get(req)
				if err != tt.results[i].err {
					t.Errorf("Unexpected error: %v", err)
				}
				if resp != tt.results[i].body {
					t.Errorf("Wrong response. Expected: %s, got: %s", tt.results[i].body, resp)
				}
			}
		})
	}
}

