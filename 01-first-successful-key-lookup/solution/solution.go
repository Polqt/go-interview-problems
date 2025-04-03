package solution

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("key not found")

func get(ctx context.Context, address string, key string) (string, error) {
	return "", nil
}

func Get(ctx context.Context, addresses []string, key string) (string, error) {
	// Creating context with timeout, so we can cancel it when receive first response (or timeout)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Channels MUST be buffered, in other case there is a goroutine leakage
	resCh, errCh := make(chan string, len(addresses)), make(chan error, len(addresses))

	for _, address := range addresses {
		go func() {
			if val, err := get(ctx, address, key); err != nil {
				errCh <- err
			} else {
				// There is a potential goroutine leak, if channel was unbuffered.
				// If the result is not first, we WILL NOT read this channel
				// and this goroutine will stuck forever
				resCh <- val
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
