package chan_queue

import (
	"errors"
)

type ChanQueue[T any] struct {
	queue    []T
	pushChan chan []T
	popChan  chan chan item[T]
}

func NewChanQueue[T any]() ChanQueue[T] {
	q := ChanQueue[T]{
		pushChan: make(chan []T),
		popChan:  make(chan chan item[T]),
	}
	go q.start()
	return q
}

// NOT THREAD SAFE
// should only be called post instantiation
func (q *ChanQueue[T]) start() {
	for {
		select {
		case items := <-q.pushChan:
			q.queue = append(q.queue, items...)

		case responseChan := <-q.popChan:
			responseChan <- q.pop()
		}
	}
}

// push a slice of items onto the queue
func (q ChanQueue[T]) Push(items ...T) {
	q.pushChan <- items
}

// pop an item from the queue.
// if queue is empty, "empty queue" error is present
func (q ChanQueue[T]) Pop() (T, error) {
	responseChan := make(chan item[T])
	q.popChan <- responseChan

	response := <-responseChan
	return response.unpack()
}

// NOT THREAD SAFE
// should only be called in start
func (q *ChanQueue[T]) pop() item[T] {
	if q.empty() {
		return newItem(*new(T), errors.New("empty queue"))
	}

	item := q.queue[0]
	q.queue = q.queue[1:]

	return newItem(item, nil)
}

// NOT THREAD SAFE
func (q *ChanQueue[T]) empty() bool {
	return len(q.queue) == 0
}
