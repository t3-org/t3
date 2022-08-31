package hecho

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestIDConfig defines the config for RequestID middleware.
type RequestIDConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Generator defines a function to generate an ID.
	// Optional. Default value random.String(32).
	Generator func() string
}

// DefaultRequestIDConfig is the default RequestID middleware config.
var DefaultRequestIDConfig = RequestIDConfig{
	Skipper:   middleware.DefaultSkipper,
	Generator: uuidGenerator,
}

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = uuidGenerator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
			}

			c.Set(ContextKeyHexaRequestID, rid)
			// Set on the request,later maybe user need to it (e.g initializing context logger)
			req.Header.Set(echo.HeaderXRequestID, rid)
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}
