package config

type stack struct {
	keys  []string
	index int
}

func newStack(keys []string) *stack {
	return &stack{
		keys:  keys,
		index: len(keys),
	}
}

func (s *stack) Push(key string) {
	s.keys = append(s.keys, key)
	s.index++
}

func (s *stack) Pop() string {
	if s.index == 0 {
		return ""
	}

	key := s.keys[s.index-1]
	s.index--
	return key
}

func (s *stack) Len() int {
	return s.index
}

func (s *stack) IsEmpty() bool {
	return s.index == 0
}
