package build

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type MockDoer struct {
	stop chan struct{}
	done chan struct{}
}

func NewMockDoer() MockDoer {
	return MockDoer{stop: make(chan struct{}), done: make(chan struct{}, 1)}
}

func (m MockDoer) reset() {
	select {
	case <-m.done:
	default:
	}
}

func (m MockDoer) Do(sleep_time int) {
	m.reset()

	select {
	case <-time.After(time.Duration(sleep_time) * time.Millisecond):
		// fmt.Printf("slept %d ms\n", sleep_time)
		m.done <- struct{}{}
	case <-m.stop:
		fmt.Println("stopped")
	}
}

func (m MockDoer) Cancel(sleep_time int) {
	select {
	case <-m.done:
		fmt.Println("already done")
	case m.stop <- struct{}{}:
	// max time to wait to try and stop Doer
	case <-time.After(time.Duration(2) * time.Second):
		fmt.Println("timed out")
	}
}

func TestWorkerShim(t *testing.T) {
	cancel := make(chan chan struct{})
	jobs := make(chan int)
	results := make(chan CompletedJob[int])

	workers := 5

	for i := 0; i < workers; i++ {
		doer := NewMockDoer()
		go WorkerShim[int](doer, jobs, results, cancel)
	}

	go func() {
		for {
			select {
			case jobs <- rand.Intn(150) + 50:
			case <-results:
			}
		}
	}()

	// simulate wait time
	time.Sleep(time.Duration(rand.Intn(300)+200) * time.Millisecond)

	// cancel the workers
	canceled := make(chan struct{})
	for i := 0; i < workers; i++ {
		c := make(chan struct{})
		cancel <- c
		go func(c chan struct{}) { canceled <- <-c }(c)
	}

	for i := 0; i < workers; i++ {
		<-canceled
	}
	select {
	case <-canceled:
		t.Fatal("canceled chan should be empty")
	default:
	}
}
