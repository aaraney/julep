package build

type Manager struct {
	jobs     chan<- Job
	progress <-chan JobStatus
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
			buffer = append(buffer, p.Children()...)
		}

		if activeJobs == 0 && len(buffer) == 0 {
			// done
			return
		}

		var idx int
	loop:
		for idx = 0; idx < len(buffer); idx++ {
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
