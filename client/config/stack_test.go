package config

import "testing"

func TestStack(t *testing.T) {
	t.Run("new stack", func(t *testing.T) {
		keys := []string{"a", "b", "c"}
		s := newStack(keys)
		if s.Len() != 3 {
			t.Errorf("newStack() length = %v, want %v", s.Len(), 3)
		}
		if s.IsEmpty() {
			t.Error("newStack() should not be empty")
		}
	})

	t.Run("push and pop", func(t *testing.T) {
		s := newStack([]string{})
		if !s.IsEmpty() {
			t.Error("newStack() should be empty")
		}

		s.Push("first")
		if s.Len() != 1 {
			t.Errorf("Push() length = %v, want %v", s.Len(), 1)
		}

		s.Push("second")
		if s.Len() != 2 {
			t.Errorf("Push() length = %v, want %v", s.Len(), 2)
		}

		val := s.Pop()
		if val != "second" {
			t.Errorf("Pop() = %v, want %v", val, "second")
		}
		if s.Len() != 1 {
			t.Errorf("Pop() length = %v, want %v", s.Len(), 1)
		}

		val = s.Pop()
		if val != "first" {
			t.Errorf("Pop() = %v, want %v", val, "first")
		}
		if s.Len() != 0 {
			t.Errorf("Pop() length = %v, want %v", s.Len(), 0)
		}
		if !s.IsEmpty() {
			t.Error("Stack should be empty after popping all items")
		}
	})

	t.Run("pop from empty stack", func(t *testing.T) {
		s := newStack([]string{})
		val := s.Pop()
		if val != "" {
			t.Errorf("Pop() from empty stack = %v, want empty string", val)
		}
	})

	t.Run("multiple pushes and pops", func(t *testing.T) {
		s := newStack([]string{"initial"})
		s.Push("one")
		s.Push("two")
		s.Push("three")

		if s.Len() != 4 {
			t.Errorf("Stack length = %v, want %v", s.Len(), 4)
		}

		expected := []string{"three", "two", "one", "initial"}
		for _, want := range expected {
			got := s.Pop()
			if got != want {
				t.Errorf("Pop() = %v, want %v", got, want)
			}
		}

		if !s.IsEmpty() {
			t.Error("Stack should be empty after popping all items")
		}
	})
}
