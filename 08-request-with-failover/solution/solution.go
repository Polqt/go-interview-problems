package main

import (
	"context"
	"errors"
	"time"
)

type Getter interface {
	Get(ctx context.Context, address string) ([]byte, error)
}

var (
	ErrRequestsFailed = errors.New("requests failed")
	getter            Getter
)

// RequestWithFailover attempts to request a data from available addresses:
// 1. If error, immediately try the address without waiting
// 2. If an address doesn't respond within 500ms, try the next but keep the original request running
// 3. Return the first successful response, or all ErrTaskFailed if all nodes fail
// 4. Properly handle context cancellation throughout the process
func RequestWithFailover(ctx context.Context, addresses []string) ([]byte, error) {
	ch := make(chan []byte, 1)
	errCh := make(chan error, len(addresses))

	var errCnt int
	for _, address := range addresses {
		nodeCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		go func() {
			resp, err := getter.Get(ctx, address)
			if err != nil {
				errCh <- err
				return
			}
			select {
			case ch <- resp:
			default:
				return
			}
		}()

		select {
		case res := <-ch:
			return res, nil
		case <-errCh:
			errCnt++
		case <-nodeCtx.Done():
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	if errCnt == len(addresses) {
		return nil, ErrRequestsFailed
	}

	for {
		select {
		case res := <-ch:
			return res, nil
		case <-errCh:
			errCnt++
			if errCnt == len(addresses) {
				return nil, ErrRequestsFailed
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
