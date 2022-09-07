package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aaraney/inlet/image"
	"github.com/aaraney/inlet/image_map"
)

func findDockerfiles(location string, dockerfiles chan<- string) {
	filepath.WalkDir(location, func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(d.Name(), "Dockerfile") {
			dockerfiles <- path
		}
		return err
	})
	close(dockerfiles)
}

func main() {
	dockerfiles := make(chan string)

	imgs := make(image_map.ImageMap)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for dockerfile := range dockerfiles {
			image_name := image.BuildName(dockerfile)
			// TODO: handle this err
			f, _ := os.Open(dockerfile)
			source_name, _ := image.ImageSource(f)
			f.Close()
			imgs.Insert(source_name.Name, image_map.DockerTermini{Path: dockerfile, Image: image.Image{Name: image_name}})
		}
		wg.Done()
	}()

	findDockerfiles(".", dockerfiles)
	wg.Wait()
	fmt.Printf("%#v\n", imgs)
	// image_map.ImageMap
	// image.Image
}
