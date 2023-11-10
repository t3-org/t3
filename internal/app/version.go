package app

import (
	"github.com/kamva/hexa"
	"t3.org/t3/infra"
)

// Version is the app version.
// Inject the version using `-ldflags`.
// go run -ldflags "-X t3.org/t3/internal/app.Version=`git rev-parse HEAD`" ./cmd/space/main.go version
var Version string

// Banner function print the app's banner
func Banner(product string) {
	hexa.PrintBanner(infra.Banner(), product, Version, "")
}
