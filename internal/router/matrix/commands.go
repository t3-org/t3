package matrix

import (
	"github.com/kamva/hexa"
	"space.org/space/internal/app"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/services"
)

func RegisterCommands(r hexa.ServiceRegistry, router *Router) {
	s := services.New(r)
	a := r.Service(registry.ServiceNameApp).(app.App)

	homeRes := newHomeResource(s, router, a)
	router.NotFoundHandler = homeRes.NotFound
	router.ErrHandler = homeRes.ErrorHandler

	registerHomeCommands(router, homeRes)
	registerTicketCommands(router, newTicketResource(s, a))
}
