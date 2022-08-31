package hecho

import (
	"fmt"

	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return RecoverWithConfig(middleware.DefaultSkipper)
}

// RecoverWithConfig returns a Recover middleware with config.
// See `Recover()` document.
func RecoverWithConfig(skipper middleware.Skipper) echo.MiddlewareFunc {
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					var ok bool
					err, ok = r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					err = tracer.Trace(err)
					//c.Error(tracer.Trace(err))
				}
			}()
			return next(c)
		}
	}
}
