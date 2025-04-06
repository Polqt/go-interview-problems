package main

import (
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

type MockConnection struct {
	delay time.Duration
	ready bool
	sync.Mutex
}

func (c *MockConnection) Connect() {
	c.Lock()
	defer c.Unlock()

	<-time.After(c.delay)
	c.ready = true
}

func (c *MockConnection) Disconnect() {
	c.Lock()
	defer c.Unlock()

	<-time.After(c.delay)
	c.ready = false
}

func (c *MockConnection) Send(req string) (string, error) {
	c.Lock()
	defer c.Unlock()

	if !c.ready {
		return "", errors.New("connection is not ready")
	}

	// Sending request
	<-time.After(c.delay)
	return "resp:" + req, nil
}

type MockCreator struct {
	delay       time.Duration
	connections []Connection
	sync.Mutex
}

func NewMockCreator(max int, delay time.Duration) *MockCreator {
	return &MockCreator{connections: make([]Connection, 0, max), delay: delay}
}

func (c *MockCreator) NewConnection() (Connection, error) {
	c.Lock()
	defer c.Unlock()

	for i, conn := range c.connections {
		if conn.(*MockConnection).ready {
			c.connections = slices.Delete(c.connections, i, i+1)
		}
	}

	if len(c.connections) == cap(c.connections) {
		return nil, errors.New("too many connections")
	}

	conn := &MockConnection{
		ready: false,
		delay: c.delay,
	}
	c.connections = append(c.connections, conn)
	return conn, nil
}

type UnsafeStorage struct {
	delay time.Duration
	sem   chan struct{}
	data  []string
	sync.Mutex
}

func NewUnsafeStorage(delay time.Duration) *UnsafeStorage {
	return &UnsafeStorage{sem: make(chan struct{}, 1), delay: delay}
}

func (s *UnsafeStorage) Save(data string) {
	select {
	case s.sem <- struct{}{}:
		<-s.sem
	default:
		data = "" // corrupt string
	}
	<-time.After(s.delay)

	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, data)
}

func TestSendAndSave(t *testing.T) {
	tests := []struct {
		name     string
		requests []string
		maxConn  int
		delay    time.Duration
		ttl      time.Duration
		err      error
	}{
		{
			name:     "multiple requests",
			requests: []string{"req1", "req2", "req3", "req4", "req5"},
			maxConn:  2,
			delay:    50 * time.Millisecond,
			ttl:      425 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saver := NewUnsafeStorage(tt.delay)
			creator := NewMockCreator(tt.maxConn, tt.delay)

			requests := make([]string, len(tt.requests))
			copy(requests, tt.requests)
			SendAndSave(creator, saver, requests, tt.maxConn)

			for _, conn := range creator.connections {
				if conn.(*MockConnection).ready {
					t.Errorf("Connection is not closed")
				}
			}

			if len(saver.data) != len(tt.requests) {
				t.Errorf("Expected %d saved items, got %d", len(saver.data), len(tt.requests))
			}

			m := map[string]bool{}
			for _, req := range tt.requests {
				m[req] = true
			}

			for _, data := range saver.data {
				if data == "" {
					t.Errorf("data is corrupted (empty string)")
				}
				if m[data] {
					m[data] = false
				}
			}
		})
	}
}
