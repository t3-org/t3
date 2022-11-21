package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "space.org/space/internal/registry/provider"
	_ "space.org/space/internal/router/api/doc"

	"space.org/space/cmd/space/commands"
)

func main() {
    if err := commands.Run(); err != nil {
        os.Exit(1)
    }
}


