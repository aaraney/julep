package utils

import (
	"testing"
)

func TestFindPredDir(t *testing.T) {
	// TODO: write tests windows compliant tests
	dir, err := FindPredDir("/", "bin")
	expected := "/bin"
	if err != nil {
		t.Fatal("/bin should exist on a unix system")
	}
	if dir != expected {
		t.Fatalf("%q should eq %q", dir, expected)

	}
}

func TestFindPredDirNegative(t *testing.T) {
	// TODO: write tests windows compliant tests
	_, err := FindPredDir("/", "/")
	if err == nil {
		t.Fatal("should have failed")
	}
}
