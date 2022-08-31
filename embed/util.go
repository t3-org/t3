package embed

import (
	"embed"

	"github.com/kamva/tracer"
)

func EntryNames(fs embed.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}

	return names, nil
}
