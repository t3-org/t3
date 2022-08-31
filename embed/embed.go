package embed

import (
	"embed"
	"path/filepath"
)

//go:embed migrations
var migrations embed.FS

func Migrations() embed.FS {
	return migrations
}

func MigrationsBaseDir(driver string) string {
	return filepath.Join("migrations", driver)
}
