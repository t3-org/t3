package main

import (
	"os"

	"space.org/space/cmd/space/commands"
	_ "space.org/space/internal/router/api/doc"
)

func main() {
	if err := commands.Run(); err != nil {
		os.Exit(1)
	}
}
