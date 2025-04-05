package main

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	val   string
	valid int64
}

type TtlCache struct {
	m      map[string]entry
	cancel context.CancelFunc
	sync.RWMutex
}

func NewTtlCache() *TtlCache {
	ctx, cancel := context.WithCancel(context.Background())
	cache := &TtlCache{
		m:      map[string]entry{},
		cancel: cancel,
	}
	go cache.clear(ctx)

	return cache
}

func (c *TtlCache) Set(key string, value string, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	var valid int64
	if ttl > 0 {
		valid = time.Now().Add(ttl).UnixNano()
	}
	c.m[key] = entry{val: value, valid: valid}
}

func (c *TtlCache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()

	now := time.Now().UnixNano()
	if e, ok := c.m[key]; ok && (e.valid == 0 || now < e.valid) {
		return e.val, true
	}

	return "", false
}

func (c *TtlCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, key)
}

func (c *TtlCache) Stop() {
	c.cancel()
}

func (c *TtlCache) clear(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Lock()

			now := time.Now().UnixNano()
			for k, e := range c.m {
				if e.valid != 0 && now > e.valid {
					delete(c.m, k)
				}
			}

			c.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
