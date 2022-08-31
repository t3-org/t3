package hredis

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type Service struct {
	*redis.Pool
}

func (h *Service) HealthIdentifier() string {
	return "redis"
}

func (h *Service) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	con := h.Pool.Get()
	if _, err := con.Do("PING"); err != nil {
		hlog.Error("error on send ping to Redis", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusDead
	}

	return hexa.StatusAlive
}

func (h *Service) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	con := h.Pool.Get()
	if _, err := con.Do("PING"); err != nil {
		hlog.Error("error on send ping to Redis", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusUnReady
	}
	return hexa.StatusReady
}

func (h *Service) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    h.HealthIdentifier(),
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}

var _ hexa.Health = &Service{}
