package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		values [][]int
		result []int
	}{
		{
			name:   "empty input",
			values: [][]int{},
			result: []int{},
		},
		{
			name:   "single empty channel",
			values: [][]int{{}},
			result: []int{},
		},
		{
			name:   "single channel with values",
			values: [][]int{{3, 1, 4}},
			result: []int{1, 3, 4},
		},
		{
			name:   "multiple channels with balanced load",
			values: [][]int{{9, 3}, {8, 2}, {7, 1}},
			result: []int{1, 2, 3, 7, 8, 9},
		},
		{
			name:   "channels with uneven distributions",
			values: [][]int{{5}, {}, {1, 2, 3, 4}},
			result: []int{1, 2, 3, 4, 5},
		},
		{
			name:   "channels with duplicate values",
			values: [][]int{{1, 3}, {1, 3}, {3}},
			result: []int{1, 1, 3, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := []<-chan int{}
			for _, chValues := range tt.values {
				ch := make(chan int)
				channels = append(channels, ch)
				go func() {
					for _, val := range chValues {
						ch <- val
					}
					close(ch)
				}()
			}
			out := Merge(channels...)

			result := []int{}
			for val := range out {
				result = append(result, val)
			}

			sort.Ints(tt.result)
			sort.Ints(result)

			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("Expected: %+q, got: %+q", tt.result, result)
			}
		})
	}
}
