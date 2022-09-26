package build

// TODO rework interface to incorporate errors
type Doer[T any] interface {
	Do(T)
	Cancel(T)
}

func WorkerShim[T any](doer Doer[T], jobs <-chan T, done_jobs chan<- CompletedJob[T], cancel chan chan struct{}) {
	done := make(chan struct{}, 1)

	for {
		select {

		case job := <-jobs:

			// start job in background
			go func(job T) {
				defer func() { done <- struct{}{} }()
				doer.Do(job)
			}(job)

			// wait for the job to be done
			// or
			// cancel signal
			select {
			case <-done:
				// TODO: handle job errors
				done_jobs <- CompletedJob[T]{Job: job, Err: nil}

			// job is in progress
			case c := <-cancel:
				// should block until job has been canceled
				doer.Cancel(job)
				// consider deadline function in future
				// Deadline(func() { doer.Cancel(job) }, time.Duration(2)*time.Second)
				c <- struct{}{}
				return
			}

		case c := <-cancel:
			c <- struct{}{}
			return
		}
	}
}

// Deadline calls fn in go routine, chan is filled with true if fn finishes before timeout. Chan is
// filled with false if timeout expires before fn returns.
// func Deadline(fn func(), timeout time.Duration) <-chan bool {
// 	done := make(chan struct{})
// 	finished := make(chan bool, 1)
// 	go func() {
// 		defer func() { done <- struct{}{} }()
// 		fn()
// 	}()

// 	select {
// 	case <-time.After(timeout):
// 		finished <- false
// 	case <-done:
// 		finished <- true
// 	}

// 	return finished
// }
