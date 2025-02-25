package queue

import (
	"sync"
	"testing"
)

// Test cases using table-driven tests
func TestQueueOperations(t *testing.T) {
	tests := []struct {
		name        string
		operations  func(q *Queue[int]) []int
		expected    []int
		expectError bool
	}{
		{
			name: "Enqueue and Dequeue",
			operations: func(q *Queue[int]) []int {
				q.Enqueue(10)
				q.Enqueue(20)
				val1, _ := q.Dequeue()
				val2, _ := q.Dequeue()
				return []int{val1, val2}
			},
			expected:    []int{10, 20},
			expectError: false,
		},
		{
			name: "Dequeue from empty queue",
			operations: func(q *Queue[int]) []int {
				val, err := q.Dequeue()
				if err != nil {
					return nil
				}
				return []int{val}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Single Element Queue",
			operations: func(q *Queue[int]) []int {
				q.Enqueue(42)
				val, _ := q.Dequeue()
				return []int{val}
			},
			expected:    []int{42},
			expectError: false,
		},
		{
			name: "Peek without removing",
			operations: func(q *Queue[int]) []int {
				q.Enqueue(100)
				val, _ := q.Peek()
				return []int{val}
			},
			expected:    []int{100},
			expectError: false,
		},
		{
			name: "Peek on empty queue",
			operations: func(q *Queue[int]) []int {
				val, err := q.Peek()
				if err != nil {
					return nil
				}
				return []int{val}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Check IsEmpty on empty queue",
			operations: func(q *Queue[int]) []int {
				if q.IsEmpty() {
					return []int{1} // Represents true
				}
				return []int{0} // Represents false
			},
			expected:    []int{1},
			expectError: false,
		},
		{
			name: "Check IsEmpty after enqueue",
			operations: func(q *Queue[int]) []int {
				q.Enqueue(1)
				if q.IsEmpty() {
					return []int{1} // Should be false
				}
				return []int{0} // Should be true
			},
			expected:    []int{0},
			expectError: false,
		},
		{
			name: "Check Len after multiple operations",
			operations: func(q *Queue[int]) []int {
				q.Enqueue(1)
				q.Enqueue(2)
				q.Dequeue()
				return []int{q.Len()}
			},
			expected:    []int{1},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New[int]()
			result := tt.operations(q)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range tt.expected {
				if result == nil || i >= len(result) || result[i] != v {
					t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
				}
			}
		})
	}
}

// Test concurrent access to the queue
func TestQueueConcurrency(t *testing.T) {
	q := New[int]()
	wg := sync.WaitGroup{}
	numOperations := 1000

	// Step 1: Concurrently enqueue items
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(i)
		}(i)
	}

	// Wait for all enqueues to complete
	wg.Wait()

	// Step 2: Concurrently dequeue items
	results := make([]int, 0, numOperations)
	mu := sync.Mutex{} // Protects results slice
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			item, err := q.Dequeue()
			if err == nil {
				mu.Lock()
				results = append(results, item)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Step 3: Ensure the queue is empty
	if q.Len() != 0 {
		t.Errorf("Expected empty queue, but got %d items remaining", q.Len())
	}
}
