package hecho

import (
	"fmt"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/htel"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	Skipper middleware.Skipper

	Propagator     propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
	ServerName     string
}

const (
	instrumentationName = "github.com/kamva/hecho"
)

// Tracing Enables distributed tracing using openTelemetry library.
// In echo if a handler panic error, it will catch by the `Recover`
// middleware, to get panic errors too, please use this middleware
// before the Recover middleware, so it will get the recovered
// errors too.
// You can use TracingDataFromUserContext middleware to set user_id
// and correlation_id too.
func Tracing(cfg TracingConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	tracer := cfg.TracerProvider.Tracer(instrumentationName)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			r := c.Request()
			opts := []trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(cfg.ServerName, c.Path(), r)...),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", r.Method)
			}

			// Extract the parent from the request, but this is a gateway that users
			// send request to it, check if propagation from external requests has any
			// security issue.
			ctx := cfg.Propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer func() { span.End() }()

			c.SetRequest(r.Clone(ctx))

			err := next(c)
			if err != nil {
				c.Error(err) // apply the error to set the response code
			}

			span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)...)

			// Set span status:
			if c.Response().Status >= 500 && err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(c.Response().Status))
			}

			return nil // we applied the error, so we don't need to return it again.
		}
	}
}

// TracingDataFromUserContext sets some tags,... on tracing span
// using hexa context. This middleware should be after hexa context
// middleware because if needs to the hexa context.
func TracingDataFromUserContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			user := hexa.CtxUser(ctx)
			if user == nil {
				return next(c)
			}

			// Add user's id, correlation_id
			span := trace.SpanFromContext(ctx)
			span.SetAttributes(
				semconv.EnduserIDKey.String(user.Identifier()),                 // enduser.id
				htel.EnduserUsernameKey.String(user.Username()),                // enduser.username
				semconv.EnduserRoleKey.String(strings.Join(user.Roles(), ",")), // enduser.role
				htel.CorrelationIDKey.String(hexa.CtxCorrelationId(ctx)),       // ctx.correlation_id
			)

			return next(c)
		}
	}
}
