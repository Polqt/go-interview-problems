package main

import (
	"context"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) == 0 {
		return "", nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Channels MUST be buffered, in other case there is a goroutine leakage
	resCh, errCh := make(chan string, 1), make(chan error, len(addresses))

	for _, address := range addresses {
		go func() {
			if val, err := getter.Get(ctx, address, key); err != nil {
				errCh <- err
			} else {
				// There is a potential goroutine leak, if channel was unbuffered.
				// If the result is not first, we WILL NOT read this channel
				// and this goroutine will stuck forever
				select {
				case resCh <- val:
				default:
				}
			}
		}()
	}

	var errCount int
	for {
		select {
		case err := <-errCh:
			// If error count is equal to addresses count
			// it means that no goroutine left and we can return an error
			errCount++
			if errCount == len(addresses) {
				return "", err
			}
		case val := <-resCh:
			return val, nil
		case <-ctx.Done():
			return "", context.Canceled
		}
	}
}
