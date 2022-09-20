package utils

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// returns hash of current git HEAD node
func GitHead(initialPath string) (string, error) {
	// look for .git dir in initialPath and it's parents
	dir, err := FindPredDir(initialPath, ".git")
	if err != nil {
		return "", err
	}

	// file, HEAD, has ref to current ref file
	head_file := filepath.Join(dir, "HEAD")
	file, err := os.ReadFile(head_file)
	if err != nil {
		return "", err
	}

	words := bytes.Split(file, []byte(" "))
	if len(words) < 2 {
		return "", errors.New("invalid .git/HEAD ref file")
	}

	// format of HEAD file:
	// ref: refs/heads/main
	ref_loc := strings.TrimSuffix(string(words[1]), "\n")

	ref_file := filepath.Join(dir, ref_loc)

	head_hash, err := os.ReadFile(ref_file)
	if err != nil {
		return "", err
	}

	return string(head_hash), nil
}
