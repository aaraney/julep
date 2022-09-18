package utils

import (
	"fmt"
	"os"
	"strings"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home := os.Getenv("HOME")
		return fmt.Sprintf("%s%s", home, path[1:])
	}
	return path
}
