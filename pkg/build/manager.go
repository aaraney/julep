package build

import "sync"

type Manager struct {
	jobs     chan Job
	progress chan CompletedJob[Job]
}

func NewManager() Manager {
	jobs := make(chan Job)
	progress := make(chan CompletedJob[Job])

	return Manager{jobs: jobs, progress: progress}
}

func (manager Manager) StartJobs(inputPaths ...Job) {
	if len(inputPaths) == 0 {
		// no jobs
		return
	}

	var buffer []Job
	var activeJobs int

	for _, job := range inputPaths {
		select {
		// add job
		case manager.jobs <- job:
			activeJobs++
		default:
			buffer = append(buffer, job)
		}
	}

	// based on progress, either add jobs or prunes
	for {
		p := <-manager.progress
		activeJobs--
		if !p.Ok() {
			// implicitly prune
			// TODO: add do some alerting here
		} else {
			buffer = append(buffer, p.Job.Children()...)
		}

		if activeJobs == 0 && len(buffer) == 0 {
			// done
			return
		}

		// fill workers with jobs
		var idx int
	loop:
		for idx = 0; idx < len(buffer); idx++ {
			// need to block until we can push first job
			// without this, it is possible, although unlikely, that the below default select case
			// is chosen instead of pushing a job from the buffer. this will result in a dead lock
			// if in this iteration of the outer while loop, we pull the last result from the
			// workers and no new work is pushed to them.
			if idx == 0 {
				manager.jobs <- buffer[idx]
				activeJobs++
				continue
			}

			select {
			case manager.jobs <- buffer[idx]:
				activeJobs++
			default:
				break loop
			}
		}

		buffer = buffer[idx:]
	}

}

func (m Manager) StartBuilders(n_builders int, builder Doer[Job]) []chan chan struct{} {
	cancelers := make([]chan chan struct{}, n_builders)

	for i := 0; i < n_builders; i++ {
		cancel := make(chan chan struct{})
		cancelers[i] = cancel

		// start builders
		go WorkerShim(builder, m.jobs, m.progress, cancel)
	}

	return cancelers
}

func (Manager) CancelBuilders(cancelers []chan chan struct{}) {
	wg := sync.WaitGroup{}
	for i, cancel := range cancelers {
		wg.Add(1)

		go func(i int, cancel chan chan struct{}) {
			defer wg.Done()
			c := make(chan struct{})
			cancel <- c
			<-c
		}(i, cancel)
	}
	wg.Wait()
}
