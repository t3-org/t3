package jobs

import (
	"github.com/kamva/hexa"
	"t3.org/t3/internal/app"
)

type Resources struct {
	app app.App
}

func newResources(_ hexa.ServiceRegistry, a app.App) *Resources {
	return &Resources{
		app: a,
	}
}
