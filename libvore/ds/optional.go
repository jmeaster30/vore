package ds

type Optional[T any] struct {
	data     T
	hasValue bool
}

func None[T any]() Optional[T] {
	result := Optional[T]{}
	result.hasValue = false
	return result
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		data:     value,
		hasValue: true,
	}
}

func OptionalEqual[T comparable](left Optional[T], right Optional[T]) bool {
	if !left.HasValue() && !right.HasValue() {
		return true
	}

	if left.HasValue() != right.HasValue() {
		return false
	}

	return left.GetValue() == right.GetValue()
}

func OptionalMap[T any, U any](value Optional[T], mappingFunction func(T) U) Optional[U] {
	if value.HasValue() {
		return Some(mappingFunction(value.GetValue()))
	}
	return None[U]()
}

func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

func (o Optional[T]) GetValue() T {
	if !o.hasValue {
		panic("Attempting to read value from empty optional :(")
	}
	return o.data
}

func (o Optional[T]) GetValueOrDefault(def T) T {
	if !o.hasValue {
		return def
	}
	return o.data
}
