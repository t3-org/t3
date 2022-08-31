package api

import (
	"github.com/labstack/echo/v4"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
)

type API struct {
	Echo  *echo.Echo
	API   *echo.Group // /api/v1
	App   app.App
	SP    base.ServiceProvider
	Guest echo.MiddlewareFunc
	Auth  echo.MiddlewareFunc
	Debug echo.MiddlewareFunc
}

func (api *API) RegisterRoutes() {
	api.registerLabRoutes()
	api.registerPlanetRoutes()
}
