# Concurrent Queue
Implement a thread-safe queue with the following properties:

- It has a fixed maximum size set at initialization.
- `Push(val int) error` adds an item to the queue. If the queue is full, it should return `ErrQueueFull`.
- `Pop() int` removes and returns first item from the queue. If the queue is empty, return -1.
- `Peek() int` returns first item from the queue. If the queue is empty, return -1.
- The queue must be safe to use from multiple goroutines simultaneously.

## Tags
`Concurrency`
