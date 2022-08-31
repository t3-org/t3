package hecho

import (
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

// DebugMiddleware make a route available just in debug mode.
func DebugMiddleware(e *echo.Echo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if !e.Debug {
				return errRouteAvailableInDebugMode
			}

			return tracer.Trace(next(ctx))
		}
	}
}
