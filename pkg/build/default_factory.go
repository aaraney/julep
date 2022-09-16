package build

import (
	"fmt"

	"github.com/aaraney/inlet/image"
	"github.com/aaraney/inlet/image_map"
)

type DefaultJobFactory struct {
	image_map image_map.ImageMapPair
	tagger    Tagger
	registry  string
}

func (f *DefaultJobFactory) JobsFromPaths(paths ...string) []DefaultJob {
	var jobs []DefaultJob

	for _, path := range paths {
		name := image.BuildName(path)

		if !f.image_map.Exists(name) {
			// TODO: should log here
			continue
		}

		// TODO: refactor this. wayyy too much coupling
		jobs = append(jobs, DefaultJob{factory: f, image: image.ImagePair{Name: name, Path: path}})
	}
	return jobs
}

func (f *DefaultJobFactory) GetChildren(img image.ImagePair) []DefaultJob {
	if !f.image_map.Exists(img.Name) {
		return []DefaultJob{}
	}

	children := f.image_map[img.Name]

	jobs := make([]DefaultJob, len(children))

	for idx, child := range children {
		jobs[idx] = DefaultJob{factory: f, image: child}
	}

	return jobs
}

func (d DefaultJobFactory) GetFullName(img image.ImagePair) string {
	return fmt.Sprintf("%s/%s", d.registry, img.Name)

}

func (d *DefaultJobFactory) GetTag(img image.ImagePair) string {
	return d.tagger.Tag(img)
}
