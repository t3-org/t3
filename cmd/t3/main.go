package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "t3.org/t3/internal/registry/provider"
	_ "t3.org/t3/internal/router/api/doc"

	"t3.org/t3/cmd/t3/commands"
)

func main() {
	if err := commands.Run(); err != nil {
		os.Exit(1)
	}
}
