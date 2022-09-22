package set

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(key T) {
	s[key] = struct{}{}
}

func (s Set[T]) Delete(key T) {
	delete(s, key)
}

func (s Set[T]) In(key T) bool {
	_, ok := s[key]
	return ok
}

func (s Set[T]) Cardinality() int {
	return len(s)
}

func (s Set[T]) AsSlice() []T {
	ts := make([]T, len(s))
	var i int
	for k := range s {
		ts[i] = k
		i++
	}
	return ts
}
