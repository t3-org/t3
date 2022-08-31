package hecho

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

// EchoService implements hexa service.
type EchoService struct {
	*echo.Echo
	Address string
}

func (s *EchoService) HealthIdentifier() string {
	return "http_server"
}

func (s *EchoService) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	return hexa.StatusAlive // TODO: do real liveness check(send a ping to the echo server).
}

func (s *EchoService) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady // TODO: do real readiness check(send a ping request to the echo server).
}

func (s *EchoService) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    s.HealthIdentifier(),
		Alive: s.LivenessStatus(ctx),
		Ready: s.ReadinessStatus(ctx),
	}
}

func (s *EchoService) Run() error {
	return tracer.Trace(s.Start(s.Address))
}

func (s *EchoService) Shutdown(ctx context.Context) error {
	return tracer.Trace(s.Echo.Shutdown(ctx))
}

var _ hexa.Health = &EchoService{}
var _ hexa.Runnable = &EchoService{}
var _ hexa.Shutdownable = &EchoService{}
