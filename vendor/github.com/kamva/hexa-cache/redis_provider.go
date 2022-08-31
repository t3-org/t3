package hcache

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type redisCacheProvider struct {
	opts *RedisOptions
}

func NewRedisCacheProvider(opts *RedisOptions) Provider {
	return &redisCacheProvider{opts: opts}
}

func (p *redisCacheProvider) Cache(name string) Cache {
	return NewRedisCache(name, p.opts)
}

func (p *redisCacheProvider) HealthIdentifier() string {
	return "redis_cache_provider"
}

func (h *redisCacheProvider) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	con := h.opts.Pool.Get()
	if _, err := con.Do("PING"); err != nil {
		hlog.Error("error on send ping to Redis", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusDead
	}

	return hexa.StatusAlive
}

func (h *redisCacheProvider) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	con := h.opts.Pool.Get()
	if _, err := con.Do("PING"); err != nil {
		hlog.Error("error on send ping to Redis", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusUnReady
	}
	return hexa.StatusReady
}

func (h *redisCacheProvider) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    h.HealthIdentifier(),
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}


var _ Provider = &redisCacheProvider{}
var _ hexa.Health = &redisCacheProvider{}
