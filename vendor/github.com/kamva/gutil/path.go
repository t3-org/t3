package gutil

import (
	"path"
	"runtime"
)

// SourcePath returns the directory containing the source code that is calling this function.
func SourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

