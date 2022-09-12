package chan_queue

import (
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewChanQueue[int]()

	assert_eq := func(value, expected int) {
		if value != expected {
			t.Errorf("%d != %d", value, expected)
		}
	}

	vals := []int{1, 2, 3}
	q.Push(vals...)

	i := 0
	for item, err := q.Pop(); err == nil; item, err = q.Pop() {
		assert_eq(item, vals[i])
		i++
	}
}

func pusher(nPushes int, queue ChanQueue[int], wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < nPushes; i++ {
		queue.Push(i)
	}
}

func popper(queue ChanQueue[int], wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()

	var n int
	for _, err := queue.Pop(); err == nil; _, err = queue.Pop() {
		n++
	}

	results <- n
}

func TestQueueConcurrent(t *testing.T) {
	q := NewChanQueue[int]()

	n_workers := 20
	n_pushes := 100

	expected_total := n_workers * n_pushes

	wg := sync.WaitGroup{}

	for i := 0; i < n_workers; i++ {
		wg.Add(1)
		go pusher(n_pushes, q, &wg)
	}

	wg.Wait()

	resultChan := make(chan int)
	sumChan := make(chan int)

	go func(resultChan <-chan int, sumChan chan<- int) {
		var sum int
		for partial_sum := range resultChan {
			sum += partial_sum
		}
		sumChan <- sum
	}(resultChan, sumChan)

	for i := 0; i < n_workers; i++ {
		wg.Add(1)
		go popper(q, &wg, resultChan)
	}

	wg.Wait()
	close(resultChan)

	total := <-sumChan

	if expected_total != total {
		t.Fatalf("%#v != %#v", expected_total, total)
	}

}

func TestQueueConcurrentInterleave(t *testing.T) {
	q := NewChanQueue[int]()
	n_workers := 20
	n_pushes := 100

	expected_total := n_workers * n_pushes

	wg := sync.WaitGroup{}
	resultChan := make(chan int)
	sumChan := make(chan int)

	for i := 0; i < n_workers; i++ {
		wg.Add(2)
		go pusher(n_pushes, q, &wg)
		go popper(q, &wg, resultChan)
	}

	go func(resultChan <-chan int, sumChan chan<- int) {
		var sum int
		for partial_sum := range resultChan {
			sum += partial_sum
		}
		sumChan <- sum
	}(resultChan, sumChan)

	wg.Wait()

	// interleave above. last popper once all pushers and poppers are done interleaving
	wg.Add(1)
	popper(q, &wg, resultChan)

	close(resultChan)

	total := <-sumChan

	if expected_total != total {
		t.Fatalf("%#v != %#v", expected_total, total)
	}

}
