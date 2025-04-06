package main

import (
	"testing"
	"time"
)

func TestCanTake(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		rps      int
		min      int
		max      int
	}{
		{
			name:     "100 RPS for 1 second",
			duration: 1 * time.Second,
			rps:      100,
			min:      95,
			max:      105, // Allow some flexibility
		},
		{
			name:     "100 RPS for 500ms",
			duration: 500 * time.Millisecond,
			rps:      100,
			min:      45, // ~50 requests expected (100 RPS * 0.5s)
			max:      55, // Allow some flexibility
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewRateLimiter(tt.rps)
			timer := time.NewTimer(tt.duration)
			var total int
			for loop := true; loop; {
				if total > tt.max {
					break
				}

				select {
				case <-timer.C:
					loop = false
				default:
				}
				if limiter.CanTake() {
					total++
				}
			}

			if total < tt.min || total > tt.max {
				t.Errorf("failed rps. expeceted between %d and %d, got: %d", tt.min, tt.max, total)
			}
		})
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		rps      int
		min      int
		max      int
	}{
		{
			name:     "100 RPS for 1 second",
			duration: 1 * time.Second,
			rps:      100,
			min:      95,
			max:      105, // Allow some flexibility
		},
		{
			name:     "100 RPS for 500ms",
			duration: 500 * time.Millisecond,
			rps:      100,
			min:      45, // ~50 requests expected (100 RPS * 0.5s)
			max:      55, // Allow some flexibility
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewRateLimiter(tt.rps)
			timer := time.NewTimer(tt.duration)
			var total int
			for loop := true; loop; {
				if total > tt.max {
					break
				}

				select {
				case <-timer.C:
					loop = false
				default:
				}
				limiter.Take()
				total++
			}

			if total < tt.min || total > tt.max {
				t.Errorf("failed rps. expeceted between %d and %d, got: %d", tt.min, tt.max, total)
			}
		})
	}
}
