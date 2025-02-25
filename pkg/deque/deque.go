package deque

import (
	"container/list"
	"fmt"
	"sync"
)

// Deque represents a thread-safe double-ended queue.
type Deque[T any] struct {
	mu    sync.Mutex
	items *list.List // Internally backed by container/list
}

// New creates a new deque.
func New[T any]() *Deque[T] {
	return &Deque[T]{items: list.New()}
}

// PushFront adds an item to the front.
func (d *Deque[T]) PushFront(value T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items.PushFront(value)
}

// PushBack adds an item to the back.
func (d *Deque[T]) PushBack(value T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items.PushBack(value)
}

// PopFront removes and returns the front item.
func (d *Deque[T]) PopFront() (T, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("deque is empty")
	}

	element := d.items.Front()
	d.items.Remove(element)
	return element.Value.(T), nil
}

// PopBack removes and returns the back item.
func (d *Deque[T]) PopBack() (T, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("deque is empty")
	}

	element := d.items.Back()
	d.items.Remove(element)
	return element.Value.(T), nil
}

// PeekFront returns the front item without removing it.
func (d *Deque[T]) PeekFront() (T, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("deque is empty")
	}

	return d.items.Front().Value.(T), nil
}

// PeekBack returns the back item without removing it.
func (d *Deque[T]) PeekBack() (T, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("deque is empty")
	}

	return d.items.Back().Value.(T), nil
}

// Len returns the number of items in the deque.
func (d *Deque[T]) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.items.Len()
}
