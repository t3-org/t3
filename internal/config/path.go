package config

import (
	"fmt"
	"strings"

	"github.com/kamva/gutil"
)

// ProjectRootPath returns root path of the project without trailing slash.
func ProjectRootPath() string {
	return gutil.SourcePath() + "/../.."
}

// AssetPath returns assets path of the project without trailing slash.
func AssetPath(basePath string, path string) string {
	return generatePath(basePath, "assets", path)
}

// ResourcePath returns path the resources.
func ResourcePath(basePath string, path string) string {
	return generatePath(basePath, "res", path)
}

// generatePath generates a new path.
func generatePath(basePath string, middle string, path string) string {
	if basePath == "" {
		basePath = fmt.Sprintf("%s/%s", ProjectRootPath(), middle)
	}
	finalPath := fmt.Sprintf("%s/%s", strings.TrimRight(basePath, "/"), strings.TrimLeft(path, "/"))
	return strings.TrimRight(finalPath, "/")
}
