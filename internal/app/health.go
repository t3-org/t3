package app

import (
	"context"
	"fmt"
	"runtime"

	"github.com/kamva/hexa"
)

func (a *appCore) HealthIdentifier() string {
	return "core"
}

func (a *appCore) LivenessStatus(_ context.Context) hexa.LivenessStatus {
	return hexa.StatusAlive
}

func (a *appCore) ReadinessStatus(_ context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady
}

func (a *appCore) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id: a.HealthIdentifier(),
		Tags: map[string]string{
			"version": Version,
			"threads": fmt.Sprint(runtime.NumGoroutine()),
		},
		Alive: a.LivenessStatus(ctx),
		Ready: a.ReadinessStatus(ctx),
	}
}

var _ hexa.Health = &appCore{}
