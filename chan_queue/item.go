package chan_queue

// tuple-like convenience type
type item[T any] struct {
	value T
	err   error
}

func newItem[T any](value T, err error) item[T] {
	return item[T]{
		value: value,
		err:   err,
	}
}

func (i *item[T]) unpack() (T, error) {
	return i.value, i.err
}
