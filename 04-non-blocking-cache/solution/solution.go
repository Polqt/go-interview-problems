package main

import (
	"sync"
)

type Client interface {
	Get(address string) (string, error)
}

type task struct {
	body  string
	err   error
	ready chan struct{}
}

type Cache struct {
	client Client
	m      map[string]*task
	sync.Mutex
}

func NewCache(client Client) *Cache {
	return &Cache{client: client, m: make(map[string]*task)}
}

func (c *Cache) Get(address string) (string, error) {
	c.Lock()
	t := c.m[address]
	if t == nil {
		t = &task{ready: make(chan struct{})}
		c.m[address] = t
		c.Unlock()

		t.body, t.err = c.client.Get(address)
		close(t.ready)
		return t.body, t.err
	} else {
		c.Unlock()

		<-t.ready
		return t.body, t.err
	}
}
