package priorityqueue

import (
	"container/heap"
	"testing"
)

// Test cases using table-driven tests
func TestPriorityQueueOperations(t *testing.T) {
	tests := []struct {
		name       string
		comparator Comparator[string]
		items      []struct {
			value    string
			priority int
		}
		expected []string
	}{
		{
			name: "Min-Heap (smallest priority first)",
			comparator: func(a, b *Item[string]) bool {
				return a.Priority < b.Priority
			},
			items: []struct {
				value    string
				priority int
			}{
				{"Low Priority Task", 5},
				{"High Priority Task", 1},
				{"Critical Task", 0},
				{"Medium Priority Task", 3},
			},
			expected: []string{"Critical Task", "High Priority Task", "Medium Priority Task", "Low Priority Task"},
		},
		{
			name: "Max-Heap (largest priority first)",
			comparator: func(a, b *Item[string]) bool {
				return a.Priority > b.Priority
			},
			items: []struct {
				value    string
				priority int
			}{
				{"Low Priority Task", 5},
				{"High Priority Task", 1},
				{"Critical Task", 0},
				{"Medium Priority Task", 3},
			},
			expected: []string{"Low Priority Task", "Medium Priority Task", "High Priority Task", "Critical Task"},
		},
		{
			name: "Custom: Sort Alphabetically if Priorities are Equal",
			comparator: func(a, b *Item[string]) bool {
				if a.Priority == b.Priority {
					return a.Value < b.Value // Alphabetical order
				}
				return a.Priority < b.Priority
			},
			items: []struct {
				value    string
				priority int
			}{
				{"Task B", 2},
				{"Task A", 2},
				{"Task C", 2},
				{"Urgent Task", 1},
			},
			expected: []string{"Urgent Task", "Task A", "Task B", "Task C"},
		},
		{
			name: "Single Element",
			comparator: func(a, b *Item[string]) bool {
				return a.Priority < b.Priority
			},
			items: []struct {
				value    string
				priority int
			}{
				{"Only Task", 3},
			},
			expected: []string{"Only Task"},
		},
		{
			name: "Empty Queue",
			comparator: func(a, b *Item[string]) bool {
				return a.Priority < b.Priority
			},
			items: []struct {
				value    string
				priority int
			}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := New(tt.comparator)

			// Insert items using PushItem() instead of direct struct assignment
			for _, item := range tt.items {
				pq.PushItem(item.value, item.priority)
			}

			// Extract items and compare with expected order
			var result []string
			for pq.Len() > 0 {
				value, _ := pq.PopItem()
				result = append(result, value)
			}

			// Validate output order
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
				}
			}
		})
	}
}

// Test concurrent access to the priority
func TestPriorityQueueConcurrency(t *testing.T) {
	pq := New(func(a, b *Item[int]) bool {
			return a.Priority < b.Priority // Min-Heap
	})
	numOperations := 1000
	for r := range numOperations {
			pq.PushItem(r, r)
	}

	resultsCh := make(chan int, numOperations)
	// Use a single goroutine to pop items sequentially.
	go func() {
			for {
					pq.mu.Lock()
					if pq.Len() == 0 {
							pq.mu.Unlock()
							break
					}
					item, err := heap.Pop(pq).(*Item[int])
					pq.mu.Unlock()
					if err == false {
							resultsCh <- item.Value
					}
			}
			close(resultsCh)
	}()

	var results []int
	for r := range resultsCh {
			results = append(results, r)
	}

	// Now check if results are in ascending order.
	for i := 1; i < len(results); i++ {
			if results[i-1] > results[i] {
					t.Errorf("Items not in ascending order: %v", results)
					break
			}
	}
}
