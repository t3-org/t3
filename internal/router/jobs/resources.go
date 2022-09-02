package jobs

import (
	"github.com/kamva/hexa"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
)

type Resources struct {
	sp  base.ServiceProvider
	app app.App
}

func newResources(_ hexa.ServiceRegistry, a app.App) *Resources {
	return &Resources{
		app: a,
	}
}
