package app

import (
	"github.com/kamva/hexa"
	"space.org/space/infra"
)

const (
	// Version is the app version
	Version = "0.1.0"
)

// Banner function print the app's banner
func Banner(product string) {
	hexa.PrintBanner(infra.Banner(), product, Version, "")
}
