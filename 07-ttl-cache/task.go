package main

import (
	"time"
)

type TtlCache struct{}

func NewTtlCache() *TtlCache {
	return &TtlCache{}
}

func (c *TtlCache) Set(key string, value string, ttl time.Duration) {
}

func (c *TtlCache) Get(key string) (string, bool) {
	return "", false
}

func (c *TtlCache) Delete(key string) {
}

func (c *TtlCache) Stop() {
}
