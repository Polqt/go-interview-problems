package main

import (
	"time"
)

type Monitor interface {
	SendRate(int)
}

type Handler struct{}

func (h *Handler) Handle() {
	// Some work
	// Calculate rate of this function
}

func (h *Handler) LogRate(monitor Monitor, d time.Duration) {
}
