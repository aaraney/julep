package image_map

import (
	"testing"

	"github.com/aaraney/julep/pkg/image"
)

func TestInsert(t *testing.T) {
	i := make(ImageMap)
	key := "source"
	value := DockerTermini{Image: image.Image{Name: "target"}}

	i.Insert(key, value)
	if i[key][0] != value {
		t.Logf("map[%q][0] should equal %#v", key, value)
	}

	i.Insert(key, value)
	if i[key][1] != value {
		t.Logf("map[%q][1] should equal %#v", key, value)
	}

	if len(i[key]) != 2 {
		t.Logf("len(map[%q]) != 2. map[%q] = %#v", key, key, value)
	}
}

func TestExists(t *testing.T) {
	i := make(ImageMap)
	key := "source"
	i[key] = []DockerTermini{{Image: image.Image{Name: "target"}}}

	if !i.Exists(key) {
		t.Logf("key, %q, should exist in image map, %#v", key, i)
	}
}
