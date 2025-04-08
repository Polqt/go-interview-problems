package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

type mockMonitor struct {
	mu         sync.Mutex
	rateValues []int
}

func newMockMonitor() *mockMonitor {
	return &mockMonitor{
		rateValues: make([]int, 0),
	}
}

func (m *mockMonitor) SendRate(rate int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateValues = append(m.rateValues, rate)
}

func (m *mockMonitor) GetRateValues() []int {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]int, len(m.rateValues))
	copy(result, m.rateValues)
	return result
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		handlerCalls   int
		tickDuration   time.Duration
		monitoringTime time.Duration
		expectedRates  []int
	}{
		{
			name:           "No calls",
			handlerCalls:   0,
			tickDuration:   100 * time.Millisecond,
			monitoringTime: 250 * time.Millisecond,
			expectedRates:  []int{0, 0},
		},
		{
			name:           "Single call",
			handlerCalls:   1,
			tickDuration:   100 * time.Millisecond,
			monitoringTime: 250 * time.Millisecond,
			expectedRates:  []int{1, 0},
		},
		{
			name:           "Multiple calls",
			handlerCalls:   1000,
			tickDuration:   100 * time.Millisecond,
			monitoringTime: 250 * time.Millisecond,
			expectedRates:  []int{1000, 0},
		},
		{
			name:           "Custom tick duration",
			handlerCalls:   5,
			tickDuration:   200 * time.Millisecond,
			monitoringTime: 450 * time.Millisecond,
			expectedRates:  []int{5, 0},
		},
		{
			name:           "Cancel stops monitoring",
			handlerCalls:   5,
			tickDuration:   100 * time.Millisecond,
			monitoringTime: 150 * time.Millisecond,
			expectedRates:  []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			monitor := newMockMonitor()

			for range tt.handlerCalls {
				go func() {
					h.Handle()
				}()
			}

			time.Sleep(tt.monitoringTime)

			rateValues := monitor.GetRateValues()
			if reflect.DeepEqual(rateValues, tt.expectedRates) {
				t.Errorf("Expected: %v, got: %v", tt.expectedRates, rateValues)
			}
		})
	}
}
