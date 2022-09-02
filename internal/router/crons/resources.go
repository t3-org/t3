package crons

import (
	"github.com/kamva/hexa"
	"space.org/space/internal/app"
)

type Resources struct {
	app app.App
}

func newResources(_ hexa.ServiceRegistry, a app.App) *Resources {
	return &Resources{
		app: a,
	}
}
