package queue

import "errors"

// NOT THREAD SAFE
type Queue[T any] struct {
	queue []T
}

func (q *Queue[T]) Push(items ...T) {
	q.queue = append(q.queue, items...)
}

func (q *Queue[T]) Pop() (T, error) {
	if q.Empty() {
		return *new(T), errors.New("queue empty")
	}

	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, nil
}

func (q *Queue[T]) Peek() (T, error) {
	if q.Empty() {
		return *new(T), errors.New("queue empty")
	}

	item := q.queue[0]
	return item, nil
}

func (q *Queue[T]) Empty() bool {
	return len(q.queue) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.queue)
}
