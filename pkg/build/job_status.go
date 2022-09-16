package build

type Result uint8

const (
	// Error is default state
	Err Result = iota
	Ok
)

type JobStatus struct {
	Job
	Result Result
}

func (p JobStatus) Ok() bool {
	return p.Result == Ok
}
