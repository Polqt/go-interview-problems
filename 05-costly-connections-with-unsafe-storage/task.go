package main

import (
	"errors"
	"sync"
	"time"
)

type Connection struct {
	ready bool
	sync.Mutex
}

func NewConnection() *Connection {
	return &Connection{}
}

func (c *Connection) Connect() {
	c.Lock()
	defer c.Unlock()

	<-time.After(2 * time.Second)
	c.ready = true
}

func (c *Connection) Disconnect() {
	c.Lock()
	defer c.Unlock()

	<-time.After(2 * time.Second)
	c.ready = false
}

func (c *Connection) Send(req string) (string, error) {
	c.Lock()
	defer c.Unlock()

	if !c.ready {
		return "", errors.New("connection is not ready")
	}

	// Sending request
	<-time.After(1 * time.Second)
	return "resp:" + req, nil
}

type UnsafeStorage struct {
	sem  chan struct{}
	data []string
	sync.Mutex
}

func NewUnsafeStorage() *UnsafeStorage {
	return &UnsafeStorage{sem: make(chan struct{}, 1)}
}

func (s *UnsafeStorage) Save(data string) {
	select {
	case s.sem <- struct{}{}:
		<-s.sem
	default:
		data = "" // corrupt string
	}
	<-time.After(1 * time.Second)

	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, data)
}

var storage = NewUnsafeStorage()

// sendAndSave should send all requests concurrently using at most `maxConn` simultaneous connections.
// Responses must be saved using UnsafeStorage. Be careful: UnsafeStorage is not safe for concurrent use.
func sendAndSave(requests []string, maxConn int) {
}
