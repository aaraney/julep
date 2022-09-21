package image_map

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

func NewImageMapPair(rootDir string) ImageMapPair {
	N_WORKERS := 4

	wg := sync.WaitGroup{}

	dockerfiles := make(chan string)
	image_pairs := make(chan image.ImagePair)
	imgs := make(ImageMapPair)

	// writes to map
	go func() {
		for pair := range image_pairs {
			imgs.InsertOne(pair)
		}
	}()

	// sends pair to map writer
	for i := 0; i < N_WORKERS; i++ {
		wg.Add(1)
		go func() {
			for dockerfile := range dockerfiles {

				f, err := os.Open(dockerfile)
				if err != nil {
					continue
				}

				pair := image.NewImagePair(dockerfile, f)

				f.Close()
				image_pairs <- pair
			}
			wg.Done()
		}()

	}

	findDockerfiles(rootDir, dockerfiles)
	wg.Wait()
	close(image_pairs)
	return imgs
}

func findDockerfiles(location string, dockerfiles chan<- string) {
	filepath.WalkDir(location, func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(d.Name(), "Dockerfile") {
			dockerfiles <- path
		}
		return err
	})
	close(dockerfiles)
}