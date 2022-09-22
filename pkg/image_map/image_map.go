package image_map

import (
	"github.com/aaraney/inlet/pkg/image"
)

// assumes that all image _names_ are unique.
// image name does not include registry name.
// so, `google/cadvisor`'s image name is `cadvisor`.
type ImageMap map[string][]DockerTermini

type DockerTermini struct {
	Path string
	image.Image
}

func (m ImageMap) Insert(key string, images ...DockerTermini) {
	if !m.Exists(key) {
		m[key] = images
		return
	}

	m[key] = append(m[key], images...)
}

func (m ImageMap) Exists(key string) bool {
	_, ok := m[key]
	return ok
}
