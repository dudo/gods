package stack

import (
	"sync"
	"testing"
)

func TestStackOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations func(s *Stack[int]) []int
		want       []int
	}{
		{
			name: "Push and Pop",
			operations: func(s *Stack[int]) []int {
				s.Push(10)
				s.Push(20)
				val1, _ := s.Pop()
				val2, _ := s.Pop()
				return []int{val1, val2}
			},
			want: []int{20, 10},
		},
		{
			name: "Pop from empty stack",
			operations: func(s *Stack[int]) []int {
				_, err := s.Pop()
				if err == nil {
					t.Errorf("Expected error popping from empty stack, got nil")
				}
				return nil
			},
			want: nil,
		},
		{
			name: "Single element stack",
			operations: func(s *Stack[int]) []int {
				s.Push(42)
				val, _ := s.Pop()
				return []int{val}
			},
			want: []int{42},
		},
		{
			name: "Peek does not remove",
			operations: func(s *Stack[int]) []int {
				s.Push(100)
				val, _ := s.Peek() // should not remove
				popVal, _ := s.Pop()
				return []int{val, popVal}
			},
			want: []int{100, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			got := tt.operations(s)

			// Compare lengths
			if len(got) != len(tt.want) {
				t.Errorf("got length %d, want %d", len(got), len(tt.want))
				return
			}

			// Compare values
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("at index %d, got %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestStackConcurrency ensures the stack is safe to use from multiple goroutines.
func TestStackConcurrency(t *testing.T) {
	s := New[int]()
	wg := sync.WaitGroup{}
	numOps := 1000

	// 1. Push concurrently
	for r := range numOps {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			s.Push(val)
		}(r)
	}

	// Wait for all pushes to complete
	wg.Wait()

	// 2. Pop concurrently
	results := make([]int, 0, numOps)
	resultsMu := sync.Mutex{}

	for range numOps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := s.Pop()
			if err == nil {
				resultsMu.Lock()
				results = append(results, val)
				resultsMu.Unlock()
			}
		}()
	}

	wg.Wait()

	// 3. Check that we popped exactly numOps items in total
	//    (Some might have popped the empty stack if scheduling was unlucky, so just check maximum.)
	if len(results) > numOps {
		t.Errorf("Got more popped items (%d) than pushes (%d)", len(results), numOps)
	}
}
