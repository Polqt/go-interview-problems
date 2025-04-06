package main

type Client interface {
	Get(address string) (string, error)
}

type Cache struct {
	client Client
	// You can add new fields if needed
}

// Don't update signature of NewCache
func NewCache(client Client) *Cache {
	// TODO: Implement
	return &Cache{client: client}
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	// TODO: Implement. Right now it doesn't cache
	return c.client.Get(address)
}
