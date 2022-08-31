package commands

import (
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"space.org/space/internal/base"
	"space.org/space/internal/registry"
)

func runProbeServer(sp base.ServiceProvider) error {
	for _, d := range registry.Registry().Descriptors() {
		if health, ok := d.Instance.(hexa.Health); ok {
			sp.HealthReporter().AddToChecks(health)
		}
	}

	return tracer.Trace(sp.ProbeServer().Run())
}
