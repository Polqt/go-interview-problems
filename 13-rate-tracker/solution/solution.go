package main

import (
	"sync/atomic"
	"time"
)

type Monitor interface {
	SendRate(int)
}

type Handler struct {
	cnt atomic.Int32
}

func (h *Handler) Handle() {
	// Some work
	// Calculate rate of this function
	h.cnt.Add(1)
}

func (h *Handler) LogRate(monitor Monitor, d time.Duration) {
	go func() {
		for range time.After(d) {
			cnt := h.cnt.Swap(0)
			monitor.SendRate(int(cnt))
		}
	}()
}
