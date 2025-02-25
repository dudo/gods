package queue

import (
	"container/list"
	"fmt"
	"sync"
)

// Queue represents a thread-safe FIFO queue.
type Queue[T any] struct {
	mu    sync.Mutex
	items *list.List // Internally backed by container/list
}

// New creates a new queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{items: list.New()}
}

// Enqueue adds an item to the back of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items.PushBack(value)
}

// Dequeue removes and returns the front item.
func (q *Queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}

	element := q.items.Front()
	q.items.Remove(element)
	return element.Value.(T), nil
}

// Peek returns the front item without removing it.
func (q *Queue[T]) Peek() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}

	return q.items.Front().Value.(T), nil
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items.Len()
}

// IsEmpty checks if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items.Len() == 0
}
