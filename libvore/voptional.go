package libvore

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
