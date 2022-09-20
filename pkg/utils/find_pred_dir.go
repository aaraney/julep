package utils

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// look for dir with name `needle` at initial path and its predecessor paths
func FindPredDir(initialPath, needle string) (string, error) {

	for {
		// files are sorted by file names
		files, err := os.ReadDir(initialPath)
		if err != nil {
			return "", err
		}

		curAbsPath, err := filepath.Abs(initialPath)
		if err != nil {
			return "", err
		}

		// binary search to find .git directory
		i := sort.Search(len(files), func(i int) bool { return files[i].Name() >= needle })

		if i < len(files) && files[i].Name() == needle {
			// needle is present at files[i]
			return filepath.Join(curAbsPath, needle), nil
		} else {
			if curAbsPath == "/" {
				return "", errors.New("fatal: git repository not found (or any of the parent directories): .git")
			}
			initialPath = "../" + initialPath
		}
	}

}
