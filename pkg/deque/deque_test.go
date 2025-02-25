package deque

import (
	"sync"
	"testing"
)

// Test cases using table-driven tests
func TestDequeOperations(t *testing.T) {
	tests := []struct {
		name        string
		operations  func(d *Deque[int]) []int
		expected    []int
		expectError bool
	}{
		{
			name: "PushFront and PopFront",
			operations: func(d *Deque[int]) []int {
				d.PushFront(10)
				d.PushFront(20)
				val1, _ := d.PopFront()
				val2, _ := d.PopFront()
				return []int{val1, val2}
			},
			expected:    []int{20, 10},
			expectError: false,
		},
		{
			name: "PushBack and PopBack",
			operations: func(d *Deque[int]) []int {
				d.PushBack(10)
				d.PushBack(20)
				val1, _ := d.PopBack()
				val2, _ := d.PopBack()
				return []int{val1, val2}
			},
			expected:    []int{20, 10},
			expectError: false,
		},
		{
			name: "PushFront and PopBack",
			operations: func(d *Deque[int]) []int {
				d.PushFront(10)
				d.PushFront(20)
				val1, _ := d.PopBack()
				val2, _ := d.PopBack()
				return []int{val1, val2}
			},
			expected:    []int{10, 20},
			expectError: false,
		},
		{
			name: "PushBack and PopFront",
			operations: func(d *Deque[int]) []int {
				d.PushBack(10)
				d.PushBack(20)
				val1, _ := d.PopFront()
				val2, _ := d.PopFront()
				return []int{val1, val2}
			},
			expected:    []int{10, 20},
			expectError: false,
		},
		{
			name: "PopFront from empty deque",
			operations: func(d *Deque[int]) []int {
				val, err := d.PopFront()
				if err != nil {
					return nil
				}
				return []int{val}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "PopBack from empty deque",
			operations: func(d *Deque[int]) []int {
				val, err := d.PopBack()
				if err != nil {
					return nil
				}
				return []int{val}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "PeekFront without removing",
			operations: func(d *Deque[int]) []int {
				d.PushBack(10)
				val, _ := d.PeekFront()
				return []int{val}
			},
			expected:    []int{10},
			expectError: false,
		},
		{
			name: "PeekBack without removing",
			operations: func(d *Deque[int]) []int {
				d.PushBack(10)
				val, _ := d.PeekBack()
				return []int{val}
			},
			expected:    []int{10},
			expectError: false,
		},
		{
			name: "Len check",
			operations: func(d *Deque[int]) []int {
				d.PushFront(10)
				d.PushFront(20)
				return []int{d.Len()}
			},
			expected:    []int{2},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New[int]()
			result := tt.operations(d)

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

// Test concurrent access to the deque
func TestDequeConcurrency(t *testing.T) {
	d := New[int]()
	wg := sync.WaitGroup{}
	numOperations := 1000

	// Step 1: PushFront from multiple goroutines
	for i := range numOperations {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.PushFront(i)
		}(i)
	}

	// Ensure all pushes finish before popping
	wg.Wait()

	// Step 2: PopFront from multiple goroutines
	for range numOperations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.PopFront()
		}()
	}

	wg.Wait()

	if d.Len() != 0 {
		t.Errorf("Deque should be empty after concurrent push and pop operations, but has %d items left", d.Len())
	}
}
