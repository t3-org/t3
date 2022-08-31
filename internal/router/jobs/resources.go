package jobs

import (
	"space.org/space/internal/app"
	"space.org/space/internal/base"
)

type Resources struct {
	sp  base.ServiceProvider
	app app.App
}

func newResources(sp base.ServiceProvider, a app.App) *Resources {
	return &Resources{
		sp:  sp,
		app: a,
	}
}
