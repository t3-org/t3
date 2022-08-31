package hecho

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricsConfig struct {
	Skipper       middleware.Skipper
	MeterProvider metric.MeterProvider
}

func Metrics(cfg MetricsConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	meter := metric.Must(cfg.MeterProvider.Meter(instrumentationName))
	requestCounter := meter.NewFloat64Counter("http_requests_total")
	requestDuration := meter.NewFloat64Histogram("http_request_duration_seconds")
	requestSize := meter.NewInt64Histogram("http_response_size_bytes")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			startTime := time.Now()
			r := c.Request()

			attrs := []attribute.KeyValue{
				attribute.String("method", r.Method),
				attribute.String("handler", c.Path()),
			}

			err := next(c)
			if err != nil {
				c.Error(err) // apply the error to set the response code
			}

			attrs = append(attrs, attribute.Int("status", c.Response().Status))

			elapsed := float64(time.Since(startTime)) / float64(time.Second)

			requestCounter.Add(r.Context(), 1, attrs...)
			requestDuration.Record(r.Context(), elapsed, attrs...)
			requestSize.Record(r.Context(), c.Response().Size, attrs...)

			return nil // we applied the error, so we don't need to return it again.
		}
	}
}
