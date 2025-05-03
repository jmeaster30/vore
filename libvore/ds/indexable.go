package ds

type Indexable[T any] interface {
	Index(int) Optional[T]
}

type DefaultIndexable[T any] struct{}

func (d DefaultIndexable[T]) Index(idx int) Optional[T] {
	return None[T]()
}

func (i DefaultIndexable[T]) Subslice(startIndex int, endIndex int) []T {
	slice := make([]T, endIndex-startIndex)
	for idx := startIndex; idx <= endIndex; idx++ {
		value := i.Index(idx)
		if !value.HasValue() {
			break
		}
		slice = append(slice, value.GetValue())
	}
	return slice
}
