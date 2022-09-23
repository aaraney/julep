package build

import (
	"fmt"

	"github.com/aaraney/julep/pkg/image"
	"github.com/aaraney/julep/pkg/image_map"
	"github.com/aaraney/julep/pkg/set"
)

type DefaultJobFactory struct {
	image_map image_map.ImageMapPair
	tagger    Tagger
	registry  string
}

func (f *DefaultJobFactory) JobsFromPaths(paths ...string) []Job {
	candidatePaths := f.candidatesPaths(paths...)
	jobs := make([]Job, len(candidatePaths))

	for i, path := range candidatePaths {
		name := image.BuildName(path)

		// TODO: refactor this. wayyy too much coupling
		jobs[i] = DefaultJob{factory: f, image: image.ImagePair{Name: name, Path: path}}
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

// modified BFS. returns paths such that none depend on one another.
// example:
//   DAG: a -> b -> c
//   candidatesPaths("c", "a") -> ["a"]
func (f *DefaultJobFactory) candidatesPaths(paths ...string) []string {
	visited := set.NewSet[string]()
	// name: path
	candidates := make(map[string]string)

	for _, candPath := range paths {
		candName := image.BuildName(candPath)

		// skip candidate name not in image_map
		if !f.image_map.Exists(candName) {
			// TODO: should log here
			continue
		}

		// if visited, `candidate` is reachable from a predecessor startNode
		if visited.In(candName) {
			continue
		}
		visited.Add(candName)
		candidates[candName] = candPath

		children := f.children(candName)

		for {
			var toVisit []string

			for _, child := range children {
				if visited.In(child) {
					// already reached
					if _, ok := candidates[child]; ok {
						delete(candidates, child)
					}
					continue
				}

				visited.Add(child)
				toVisit = append(toVisit, f.children(child)...)
			}
			if len(toVisit) == 0 {
				break
			}
			children = toVisit
		}

	}

	candPaths := make([]string, len(candidates))

	var i int
	for _, v := range candidates {
		candPaths[i] = v
		i++
	}

	return candPaths
}
