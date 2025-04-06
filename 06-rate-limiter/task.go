package main

type RateLimiter struct{}

func NewRateLimiter(n int) *RateLimiter {
	return &RateLimiter{}
}

func (r *RateLimiter) CanTake() bool {
	return false
}

func (r *RateLimiter) Take() {
}
