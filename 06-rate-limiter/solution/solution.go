package solution

import "time"

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(n int) *RateLimiter {
	limit := time.Second / time.Duration(n)
	return &RateLimiter{ticker: time.NewTicker(limit)}
}

func (r *RateLimiter) CanTake() bool {
	select {
	case <-r.ticker.C:
		return true
	default:
		return false
	}
}

func (r *RateLimiter) Take() {
	<-r.ticker.C
}
