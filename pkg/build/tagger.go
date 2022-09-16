package build

import "github.com/aaraney/inlet/image"

type Tagger interface {
	Tag(image.ImagePair) string
}

type DefaultTagger struct{}

func (DefaultTagger) Tag(img image.ImagePair) string {
	// TODO: implement this
	return "0.0.1"
}
