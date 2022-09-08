package config

import (
	"path"

	"github.com/kamva/gutil"
)

// appRootPath returns root path of the app without trailing slash.
func appRootPath() string {
	return path.Join(gutil.SourcePath(), "../..")
}
