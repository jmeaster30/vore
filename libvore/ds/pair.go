package ds

type Pair[T any, U any] struct {
	left  T
	right U
}

func NewPair[T any, U any](left T, right U) Pair[T, U] {
	return Pair[T, U]{left, right}
}

func (p Pair[T, U]) Flip() Pair[U, T] {
	return NewPair(p.right, p.left)
}

func (p Pair[T, U]) Left() T {
	return p.left
}

func (p Pair[T, U]) Right() U {
	return p.right
}

func (p Pair[T, U]) Values() (T, U) {
	return p.left, p.right
}
