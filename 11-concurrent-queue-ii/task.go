package main

import (
	"errors"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct{}

func NewQueue(size int) *Queue {
	return &Queue{}
}

func (q *Queue) Push(val int) error {
	return nil
}

func (q *Queue) Pop() int {
	return -1
}

func (q *Queue) Peek() int {
	return -1
}
