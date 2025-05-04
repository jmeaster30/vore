package ds

type Indexable[T any] interface {
	Index(int) Optional[T]
}

type DefaultIndexable[T any] struct{}

func Subslice[T any](i Indexable[T], startIndex int, endIndex int) []T {
	slice := make([]T, 0, endIndex-startIndex+1)
	for idx := startIndex; idx <= endIndex; idx++ {
		value := i.Index(idx)
		if !value.HasValue() {
			break
		}
		slice = append(slice, value.GetValue())
	}
	return slice
}
