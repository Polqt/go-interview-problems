package solution

import "time"

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(n int) *RateLimiter {
	return &RateLimiter{}
}

func (r *RateLimiter) CanTake() bool {
	return false
}

func (r *RateLimiter) Take() {
}
