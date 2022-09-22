package image

import "io"

const Prefix = "Dockerfile."

type Image struct {
	Name     string
	Tag      string
	Hash     string
	Registry string
}

type ImagePair struct {
	SourceImage Image
	Name        string
	Path        string
}

func NewImagePair(image_path string, image_reader io.Reader) ImagePair {
	source, _ := ImageSource(image_reader)
	return ImagePair{
		Name:        BuildName(image_path),
		Path:        image_path,
		SourceImage: source,
	}
}
