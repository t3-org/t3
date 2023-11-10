package api

import (
	"github.com/labstack/echo/v4"
	"t3.org/t3/internal/app"
)

type API struct {
	echo        *echo.Echo
	app         app.App
	middlewares *Middlewares
	v1          *echo.Group
}

func New(e *echo.Echo, app app.App, m *Middlewares) *API {
	return &API{
		echo:        e,
		app:         app,
		middlewares: m,
	}
}

func (api *API) RegisterRoutes() {
	api.v1 = api.echo.Group("/api/v1")
	v1 := api.v1.Group

	api.registerLabRoutes(v1("/lab"), &labResource{app: api.app, e: api.echo})
	api.registerTicketRoutes(v1("/tickets"), &ticketResource{app: api.app})
	api.registerWebhookRoutes(v1("/webhooks"), &webhookResource{app: api.app})
}
