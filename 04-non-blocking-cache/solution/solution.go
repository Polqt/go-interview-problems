package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func getBody(address string) ([]byte, error) {
	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type task struct {
	body  []byte
	err   error
	ready chan struct{}
}

type Cache struct {
	m map[string]*task
	sync.Mutex
}

func NewCache() *Cache {
	return &Cache{m: make(map[string]*task)}
}

func (c *Cache) Get(address string) ([]byte, error) {
	c.Lock()
	t := c.m[address]
	if t == nil {
		t = &task{ready: make(chan struct{})}
		c.m[address] = t
		c.Unlock()

		t.body, t.err = getBody(address)
		close(t.ready)
		return t.body, t.err
	} else {
		c.Unlock()

		<-t.ready
		return t.body, t.err
	}
}

func main() {
	c := NewCache()
	body, err := c.Get("https://www.google.com/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
