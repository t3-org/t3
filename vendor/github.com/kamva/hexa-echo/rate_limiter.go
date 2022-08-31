package hecho

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
)

var ErrTooManyRequests = hexa.NewError(http.StatusTooManyRequests, "lib.http.too_many_requests_error")

type KeyExtractor func(c echo.Context) (string, error)

type RateLimiterConfig struct {
	Skipper      middleware.Skipper
	RateLimiter  limiter.Store
	KeyExtractor KeyExtractor
}

func RateLimiterByIP(rl limiter.Store) echo.MiddlewareFunc {
	return RateLimiter(RateLimiterConfig{
		RateLimiter:  rl,
		KeyExtractor: IPKeyExtractor,
	})
}

func RateLimiter(cfg RateLimiterConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key, err := cfg.KeyExtractor(c)
			if err != nil {
				return tracer.Trace(err)
			}

			limit, remaining, reset, ok, err := cfg.RateLimiter.Take(c.Request().Context(), key)
			if err != nil {
				hlog.Error("error on checking rate limit", hlog.Err(tracer.Trace(err)))
				return tracer.Trace(err)
			}

			resetTime := time.Unix(0, int64(reset)).UTC().Format(time.RFC1123)
			c.Response().Header().Set(httplimit.HeaderRateLimitLimit, strconv.FormatUint(limit, 10))
			c.Response().Header().Set(httplimit.HeaderRateLimitRemaining, strconv.FormatUint(remaining, 10))
			c.Response().Header().Set(httplimit.HeaderRateLimitReset, resetTime)

			// Fail if there were no tokens remaining.
			if !ok {
				c.Response().Header().Set(httplimit.HeaderRetryAfter, resetTime)
				return ErrTooManyRequests
			}

			return next(c)
		}
	}
}

func IPKeyExtractor(c echo.Context) (string, error) {
	return c.RealIP(), nil
}
