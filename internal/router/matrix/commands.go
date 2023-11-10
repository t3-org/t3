package matrix

import (
	"github.com/kamva/hexa"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
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
