package api

import (
	hecho "github.com/kamva/hexa-echo"
	"github.com/labstack/echo/v4"
)

type Middlewares struct {
	Guest echo.MiddlewareFunc
	Auth  echo.MiddlewareFunc
	Debug echo.MiddlewareFunc
}

func NewMiddlewares(e *echo.Echo) *Middlewares {
	return &Middlewares{
		Guest: hecho.GuestMiddleware(),
		Auth:  hecho.AuthMiddleware(),
		Debug: hecho.DebugMiddleware(e),
	}
}
