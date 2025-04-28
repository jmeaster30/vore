package ds

type Stack[T any] struct {
	store []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Copy() *Stack[T] {
	result := NewStack[T]()

	for _, value := range s.store {
		result.Push(value)
	}

	return result
}

// peek
func (s *Stack[T]) Peek() Optional[T] {
	if s.IsEmpty() {
		return None[T]()
	}
	return Some(s.store[len(s.store)-1])
}

// push
func (s *Stack[T]) Push(value T) {
	s.store = append(s.store, value)
}

// pop
func (s *Stack[T]) Pop() Optional[T] {
	if s.IsEmpty() {
		return None[T]()
	}
	result := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return Some(result)
}

func (s *Stack[T]) Index(index int) Optional[T] {
	if s.IsEmpty() || index < 0 || index >= len(s.store) {
		return None[T]()
	}
	return Some(s.store[index])
}

// isEmpty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.store) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.store)
}
