package main

import (
	"context"
	"io"
	"net/http"
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

type Cache struct{}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Get(ctx context.Context, address string) ([]byte, error) {
	// Implement non-blocking cache HERE
	return getBody(address)
}
