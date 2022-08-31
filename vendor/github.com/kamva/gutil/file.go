package gutil

import (
	"io/ioutil"
	"os"
)

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(f string) bool {
	info, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadAll read file content and return it as byte slice.
func ReadAll(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	return src, err
}
