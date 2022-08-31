package hecho

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// CorrelationIDConfig defines the config for CorrelationID middleware.
	CorrelationIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value unique uuid.
		Generator func() string
	}
)

var (
	// HeaderCorrelationID is the http X-Correlation-ID header name.
	HeaderCorrelationID = "X-Correlation-ID"

	// DefaultCorrelationIDConfig is the default CorrelationID middleware config.
	DefaultCorrelationIDConfig = CorrelationIDConfig{
		Skipper:   middleware.DefaultSkipper,
		Generator: uuidGenerator,
	}
)

// CorrelationID returns a X-Request-ID middleware.
func CorrelationID() echo.MiddlewareFunc {
	return CorrelationIDWithConfig(DefaultCorrelationIDConfig)
}

// CorrelationIDWithConfig returns a X-Correlation-ID middleware with config.
func CorrelationIDWithConfig(config CorrelationIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCorrelationIDConfig.Skipper
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
			cid := req.Header.Get(HeaderCorrelationID)
			if cid == "" {
				cid = config.Generator()
			}

			c.Set(ContextKeyHexaCorrelationID, cid)
			res.Header().Set(HeaderCorrelationID, cid)

			return next(c)
		}
	}
}
