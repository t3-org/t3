package hecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LimitBodySize limits the request body size.
// n less than or equal to zero means unlimited.
// n is in bytes.
func LimitBodySize(n int64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if n > 0 {
				r := c.Request()
				r.Body = http.MaxBytesReader(c.Response(), r.Body, n)
			}

			return next(c)
		}
	}
}
