package build

type CompletedJob[T any] struct {
	Job T
	Err error
}

func (p CompletedJob[T]) Ok() bool {
	return p.Err == nil
}
