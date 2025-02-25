package priorityqueue

import (
	"container/heap"
	"fmt"
	"sync"
)

// Item represents an element in the priority queue.
type Item[T any] struct {
	Value    T
	Priority int
	index    int // Internal index for heap tracking
}

// Comparator function type (returns true if `a` should come before `b`).
type Comparator[T any] func(a, b *Item[T]) bool

// PriorityQueue is a thread-safe priority queue.
type PriorityQueue[T any] struct {
	items      []*Item[T]
	comparator Comparator[T]
	mu         sync.Mutex
}

// New creates a new concurrent priority queue with a custom comparator.
func New[T any](comparator Comparator[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		items:      []*Item[T]{},
		comparator: comparator,
	}
	heap.Init(pq)
	return pq
}

// PeekItem safely returns the highest priority item without removing it.
func (pq *PriorityQueue[T]) PeekItem() (T, error) {
	pq.mu.Lock() // Change from RLock() to Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		var zero T
		return zero, fmt.Errorf("priority queue is empty")
	}

	return pq.items[0].Value, nil
}

// PushItem safely adds a new item.
func (pq *PriorityQueue[T]) PushItem(value T, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item := &Item[T]{Value: value, Priority: priority}
	heap.Push(pq, item)
}

// PopItem safely removes the highest-priority item.
func (pq *PriorityQueue[T]) PopItem() (T, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		var zero T
		return zero, fmt.Errorf("priority queue is empty")
	}

	item := heap.Pop(pq).(*Item[T])
	return item.Value, nil
}

// Swap implements heap.Interface, swapping two elements.
func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// Less implements heap.Interface, comparing two elements using the comparator.
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.comparator(pq.items[i], pq.items[j])
}

// Push implements heap.Interface, adding an item to the heap.
func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
}

// Pop implements heap.Interface, removing and returning the highest-priority item.
func (pq *PriorityQueue[T]) Pop() any {
	n := len(pq.items)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	return item
}

// Len returns the number of items in the priority queue.
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}
