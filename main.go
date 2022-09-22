package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aaraney/julep/pkg/image"
	"github.com/aaraney/julep/pkg/image_map"
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

func buildGraph(deps image_map.ImageMapPair, image_name string) map[string]any {
	m := make(map[string]any)

	if !deps.Exists(image_name) {
		return m
	}

	for _, image := range deps[image_name] {
		// only build images with unpinned source images (i.e. latest or "")
		if len(image.SourceImage.Tag) != 0 && image.SourceImage.Tag != "latest" {
			continue
		}
		m[image.Path] = buildGraph(deps, image.Name)
	}

	return m
}

func buildGraphs(deps image_map.ImageMapPair, image_path ...string) map[string]any {
	m := make(map[string]any)

	for _, path := range image_path {
		name := image.BuildName(path)

		if !deps.Exists(name) {
			continue
		}

		m[path] = buildGraph(deps, name)
	}

	return m
}

// func main() {
// 	dockerfiles := make(chan string)

// 	imgs := make(image_map.ImageMap)

// 	wg := sync.WaitGroup{}
// 	wg.Add(1)
// 	// go func() {
// 	// 	for dockerfile := range dockerfiles {
// 	// 		image_name := image.BuildName(dockerfile)
// 	// 		// TODO: handle this err
// 	// 		f, _ := os.Open(dockerfile)
// 	// 		source_name, _ := image.ImageSource(f)
// 	// 		f.Close()
// 	// 		imgs.Insert(source_name.Name, image_map.DockerTermini{Path: dockerfile, Image: image.Image{Name: image_name}})
// 	// 	}
// 	// 	wg.Done()
// 	// }()

// 	go func() {
// 		for dockerfile := range dockerfiles {
// 			image_name := image.BuildName(dockerfile)
// 			// TODO: handle this err
// 			f, _ := os.Open(dockerfile)
// 			source_name, _ := image.ImageSource(f)
// 			f.Close()
// 			imgs.Insert(source_name.Name, image_map.DockerTermini{Path: dockerfile, Image: image.Image{Name: image_name}})
// 		}
// 		wg.Done()
// 	}()

// 	findDockerfiles(".", dockerfiles)
// 	wg.Wait()
// 	i, _ := json.MarshalIndent(&imgs, "", "    ")
// 	fmt.Printf("%s\n", i)
// 	// graph := buildGraphs(imgs, "./docker/Dockerfile.source")
// 	// g, _ := json.MarshalIndent(graph, "", "    ")
// 	// fmt.Printf("%s\n", g)

// }

func main() {
	dockerfiles := make(chan string)

	imgs := make(image_map.ImageMapPair)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for dockerfile := range dockerfiles {
			// image_name := image.BuildName(dockerfile)

			// TODO: handle this err
			f, _ := os.Open(dockerfile)
			// source_name, _ := image.ImageSource(f)

			pair := image.NewImagePair(dockerfile, f)

			f.Close()
			imgs.InsertOne(pair)
			// imgs.Insert(pair.Source.Name, image_map.DockerTermini{Path: dockerfile, Image: image.Image{Name: image_name}})
		}
		wg.Done()
	}()

	findDockerfiles(".", dockerfiles)
	wg.Wait()
	// i, _ := json.MarshalIndent(&imgs, "", "    ")
	// fmt.Printf("%s\n", i)

	graph := buildGraphs(imgs, "./docker/Dockerfile.source")
	g, _ := json.MarshalIndent(graph, "", "    ")
	fmt.Printf("%s\n", g)

}
