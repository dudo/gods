package stack

import (
	"container/list"
	"fmt"
	"sync"
)

// Stack represents a thread-safe LIFO stack.
type Stack[T any] struct {
	mu    sync.Mutex
	items *list.List // Internally backed by container/list
}

// New creates a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{items: list.New()}
}

// Push adds an item to the stack.
func (s *Stack[T]) Push(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items.PushBack(value)
}

// Pop removes and returns the last item.
func (s *Stack[T]) Pop() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}

	elem := s.items.Back()
	s.items.Remove(elem)
	return elem.Value.(T), nil
}

// Peek returns the last item without removing it.
func (s *Stack[T]) Peek() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}

	return s.items.Back().Value.(T), nil
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items.Len()
}
