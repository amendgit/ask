package X

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Lines(s string) []string {
	return strings.Split(s, "\n")
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SourceDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
