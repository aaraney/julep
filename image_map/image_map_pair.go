package image_map

import (
	"github.com/aaraney/inlet/image"
)

type ImageMapPair map[string][]image.ImagePair

func (m ImageMapPair) InsertOne(pair image.ImagePair) {
	key := pair.SourceImage.Name
	if !m.Exists(key) {
		m[key] = []image.ImagePair{pair}
		return
	}

	m[key] = append(m[key], pair)
}

func (m ImageMapPair) Insert(key string, images ...image.ImagePair) {
	if !m.Exists(key) {
		m[key] = images
		return
	}

	m[key] = append(m[key], images...)
}

func (m ImageMapPair) Exists(key string) bool {
	_, ok := m[key]
	return ok
}
