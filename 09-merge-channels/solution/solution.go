package main

import "sync"

func Merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, in := range channels {
		wg.Add(1)
		go func() {
			for val := range in {
				out <- val
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
