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

func (f *DefaultJobFactory) JobsFromPaths(paths ...string) []Job {
	var jobs []Job

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

func (f *DefaultJobFactory) GetChildren(img image.ImagePair) []Job {
	if !f.image_map.Exists(img.Name) {
		return []Job{}
	}

	children := f.image_map[img.Name]

	jobs := make([]Job, len(children))

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

func (f *DefaultJobFactory) children(key string) []string {
	children := f.image_map.Children(key)
	names := make([]string, len(children))

	for i := 0; i < len(children); i++ {
		names[i] = children[i].Name
	}
	return names
}

