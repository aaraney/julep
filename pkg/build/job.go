package build

import "github.com/aaraney/inlet/image"

type Job interface {
	// path to dockerfile
	Path() string
	// registry/image
	FullName() string
	// image tag
	Tag() string
	// Jobs based off _this_ build job
	Children() []Job
}

type DefaultJob struct {
	factory *DefaultJobFactory
	image   image.ImagePair
}

func (j DefaultJob) Path() string {
	return j.image.Path
}

func (j DefaultJob) FullName() string {
	return j.factory.GetFullName(j.image)
}

func (j DefaultJob) Tag() string {
	return j.factory.GetTag(j.image)
}

func (j DefaultJob) Children() []DefaultJob {
	return j.factory.GetChildren(j.image)
}
