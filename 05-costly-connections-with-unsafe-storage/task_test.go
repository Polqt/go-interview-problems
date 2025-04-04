package main

import (
	"testing"
	"time"
)

func TestSendAndSave(t *testing.T) {
	requests := []string{"req1", "req2", "req3", "req4", "req5"}
	maxConn := 2
	expecTime := 8 * time.Second

	start := time.Now()
	sendAndSave(requests, maxConn)
	execTime := time.Since(start).Round(time.Second)

	if len(storage.data) != len(requests) {
		t.Errorf("Expected %d saved items, got %d", len(requests), len(storage.data))
	}

	for i, data := range storage.data {
		if data == "" {
			t.Errorf("data at index %d is corrupted (empty string)", i)
		}
	}

	if execTime > expecTime {
		t.Errorf("func takes too long expected: %d seconds, got %d seconds", expecTime, execTime)
	}
}
