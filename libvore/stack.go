package libvore

type Stack[T any] struct {
	store []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// peek
func (s *Stack[T]) Peek() *T {
	if s.IsEmpty() {
		return nil
	}
	return &s.store[len(s.store)-1]
}

// push
func (s *Stack[T]) Push(value T) {
	s.store = append(s.store, value)
}

// pop
func (s *Stack[T]) Pop() *T {
	if s.IsEmpty() {
		return nil
	}
	result := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return &result
}

// isEmpty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.store) == 0
}

func (s *Stack[T]) Size() uint64 {
	return uint64(len(s.store))
}
