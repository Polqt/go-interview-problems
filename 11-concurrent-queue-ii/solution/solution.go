package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	data []int
	size int
	sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: []int{},
		size: size,
	}
}

func (q *Queue) Push(val int) error {
	q.Lock()
	defer q.Unlock()
	if len(q.data) >= q.size {
		return ErrQueueFull
	} else {
		q.data = append(q.data, val)
		return nil
	}
}

func (q *Queue) Pop() int {
	q.Lock()
	defer q.Unlock()

	if len(q.data) > 0 {
		val := q.data[0]
		q.data = q.data[1:]
		return val
	} else {
		return -1
	}
}

func (q *Queue) Peek() int {
	q.Lock()
	defer q.Unlock()

	if len(q.data) > 0 {
		return q.data[0]
	} else {
		return -1
	}
}
