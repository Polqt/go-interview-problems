package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	data []int
	w    int
	r    int
	sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: make([]int, size),
	}
}

func (q *Queue) Push(val int) error {
	q.Lock()
	defer q.Unlock()

	if q.w == q.r || q.w%len(q.data) != q.r%len(q.data) {
		q.data[q.w%len(q.data)] = val
		q.w++
		return nil
	} else {
		return ErrQueueFull
	}
}

func (q *Queue) Pop() int {
	q.Lock()
	defer q.Unlock()

	if q.r != q.w {
		val := q.data[q.r%len(q.data)]
		q.r++
		return val
	} else {
		return -1
	}
}

func (q *Queue) Peek() int {
	q.Lock()
	defer q.Unlock()

	if q.r != q.w {
		return q.data[q.r%len(q.data)]
	} else {
		return -1
	}
}
