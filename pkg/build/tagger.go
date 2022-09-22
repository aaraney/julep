package build

import (
	"github.com/aaraney/julep/pkg/image"
	"github.com/aaraney/julep/pkg/utils"
)

type Tagger interface {
	Tag(image.ImagePair) string
}

type DefaultTagger struct{}

var (
	tag    string
	tagSet bool
)

func (DefaultTagger) Tag(img image.ImagePair) string {
	if tagSet {
		return tag
	}

	// NOTE: search location will likely be provided by config in future
	hash, err := utils.GitHead(".")
	if err != nil {
		return tag
	}
	tag = hash
	tagSet = true
	return tag
}
