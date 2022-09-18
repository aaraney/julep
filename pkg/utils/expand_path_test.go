package utils

import (
	"os"
	"testing"
)

func TestExpandPath(t *testing.T) {
	expected := "/user/test/Downloads"
	input := "~/Downloads"

	original_home := os.Getenv("HOME")
	os.Setenv("HOME", "/user/test")
	dwn := ExpandPath(input)
	if expected != dwn {
		t.Fatalf("%q != %q", dwn, expected)
	}

	// tear down
	os.Setenv("HOME", original_home)
}
