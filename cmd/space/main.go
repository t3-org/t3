package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"space.org/space/cmd/space/commands"
	_ "space.org/space/internal/router/api/doc"
)

func main() {
	if err := commands.Run(); err != nil {
		os.Exit(1)
	}
}
