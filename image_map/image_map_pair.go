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
	dockerfiles := make(chan string)

	imgs := make(ImageMapPair)
	wg := sync.WaitGroup{}

	N_WORKERS := 4

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
				imgs.InsertOne(pair)
			}
			wg.Done()
		}()

	}

	findDockerfiles(rootDir, dockerfiles)
	wg.Wait()
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
