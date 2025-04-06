package main

import (
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var (
	m  = map[string]bool{}
	mu sync.RWMutex
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) ([]string, error) {
	if depth <= 0 {
		return nil, nil
	}
	mu.RLock()
	if m[url] {
		mu.RUnlock()
		return nil, nil
	}
	mu.RUnlock()

	mu.Lock()
	m[url] = true
	mu.Unlock()

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	result := []string{body}

	ch := make(chan string)
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if res, err := Crawl(u, depth-1, fetcher); err == nil {
				for _, b := range res {
					ch <- b
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for b := range ch {
		result = append(result, b)
	}

	return result, nil
}
