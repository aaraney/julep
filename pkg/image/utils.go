package image

import (
	"bufio"
	"errors"
	"io"
	"path"
	"strings"
)

// first is given precedence
// 1. `docker/Dockerfile.app` -> app
// 2. `docker/app/Dockerfile` -> docker-app
func BuildName(p string) string {
	if p == "" {
		return p
	}
	dirs, basename := path.Split(p)

	// 1.
	if strings.HasPrefix(basename, Prefix) {
		return basename[len(Prefix):]
	}

	// 2.
	if dirs != "" {
		// ignore trailing slash
		return strings.ReplaceAll(dirs[:len(dirs)-1], "/", "-")
	}

	return ""
}

// https://docs.docker.com/engine/reference/builder/#from
//
// patterns:
// FROM [--platform=<platform>] <image> [AS <name>]
// FROM [--platform=<platform>] <image>[:<tag>] [AS <name>]
// FROM [--platform=<platform>] <image>[@<digest>] [AS <name>]
//
// FROM {1} [--platform=<platform>] {2} <image>[@<digest>] {3} [AS <name>]
func parseImageSource(source string) string {
	if !strings.HasPrefix(source, "FROM") {
		return ""
	}

	// image name in either in index 1 or 2
	split := strings.Split(source, " ")

	switch len(split) {
	case 1:
		return ""
	case 2:
		return split[1]
	// greater than 2
	default:
		if strings.HasPrefix(split[1], "--platform") {
			return split[2]
		}
		return split[1]
	}

}

// <image>
// <image>[:<tag>]
// <image>[@<digest>]
// registry/etc/<image>
// registry/etc/<image>[:<tag>]
// registry/etc/<image>[@<digest>]
func ParseImageStr(image string) Image {
	var registry string

	if i := strings.LastIndex(image, "/"); i > 0 && i < len(image)-1 {
		registry = image[:i]
		image = image[i+1:]
	}

	if i := strings.IndexRune(image, ':'); i > 0 && i < len(image)-1 {
		name := image[:i]
		tag := image[i+1:]
		return Image{Name: name, Tag: tag, Registry: registry}
	}

	if i := strings.IndexRune(image, '@'); i > 0 && i < len(image)-1 {
		name := image[:i]
		hash := image[i+1:]
		return Image{Name: name, Hash: hash, Registry: registry}
	}

	return Image{Name: image, Registry: registry}
}

func ImageSource(r io.Reader) (Image, error) {
	scanner := bufio.NewScanner(r)

	var imageFullName string

	// look for first line that starts with FROM
	for scanner.Scan() {
		imageFullName = parseImageSource(scanner.Text())
		if imageFullName != "" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return Image{}, err
	}

	if imageFullName == "" {
		return Image{}, errors.New("no image source found")
	}

	return ParseImageStr(imageFullName), nil
}
