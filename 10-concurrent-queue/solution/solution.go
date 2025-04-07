package main

import "errors"

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	ch chan int
}

func NewQueue(size int) *Queue {
	return &Queue{ch: make(chan int, size)}
}

func (q *Queue) Push(val int) error {
	select {
	case q.ch <- val:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *Queue) Pop() int {
	select {
	case val := <-q.ch:
		return val
	default:
		return -1
	}
}
