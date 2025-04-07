package main

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		ops      []string
		args     []int
		expected []any
	}{
		{
			name:     "basic operations",
			size:     3,
			ops:      []string{"push", "push", "push", "pop", "pop", "pop"},
			args:     []int{1, 2, 3, 0, 0, 0},
			expected: []any{nil, nil, nil, 1, 2, 3},
		},
		{
			name:     "empty queue pop",
			size:     1,
			ops:      []string{"pop", "push", "pop"},
			args:     []int{0, 5, 0},
			expected: []any{-1, nil, 5},
		},
		{
			name:     "queue full error",
			size:     2,
			ops:      []string{"push", "push", "push", "pop", "push", "pop"},
			args:     []int{5, 10, 15, 0, 20, 0},
			expected: []any{nil, nil, ErrQueueFull, 5, nil, 10},
		},
		{
			name:     "interleaved operations",
			size:     2,
			ops:      []string{"push", "pop", "push", "push", "pop", "pop"},
			args:     []int{1, 0, 2, 3, 0, 0},
			expected: []any{nil, 1, nil, nil, 2, 3},
		},
		{
			name:     "peek operations",
			size:     3,
			ops:      []string{"push", "peek", "push", "peek", "pop", "peek", "pop", "peek"},
			args:     []int{10, 0, 20, 0, 0, 0, 0, 0},
			expected: []any{nil, 10, nil, 10, 10, 20, 20, -1},
		},
		{
			name:     "peek empty queue",
			size:     2,
			ops:      []string{"peek", "push", "peek", "pop", "peek"},
			args:     []int{0, 42, 0, 0, 0},
			expected: []any{-1, nil, 42, 42, -1},
		},
		{
			name:     "mixed operations",
			size:     5,
			ops:      []string{"push", "push", "peek", "pop", "peek", "push", "push", "push", "peek"},
			args:     []int{5, 10, 0, 0, 0, 15, 20, 25, 0},
			expected: []any{nil, nil, 5, 5, 10, nil, nil, nil, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQueue(tt.size)
			results := make([]any, len(tt.ops))

			for i, op := range tt.ops {
				var result any

				switch op {
				case "push":
					result = q.Push(tt.args[i])
				case "peek":
					result = q.Peek()
				case "pop":
					result = q.Pop()
				}

				results[i] = result
			}

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("Expected: %v, got: %v", tt.expected, results)
			}
		})
	}
}

func TestConcurrentQueue(t *testing.T) {
	queue := NewQueue(1000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := range 1000 {
		go func() {
			defer wg.Done()
			queue.Push(i + 1)
		}()
	}

	ch := make(chan int, 1000)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if val := queue.Pop(); val > -1 {
					ch <- val
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	var result []int
	for val := range ch {
		result = append(result, val)
	}

	expected := make([]int, 1000)
	for i := range 1000 {
		expected[i] = i + 1
	}

	sort.Ints(result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: 1..1000, got: %v", result)
	}
}
