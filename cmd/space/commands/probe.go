package commands

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/probe"
	"github.com/kamva/tracer"
	"space.org/space/internal/registry"
)

func runProbeServer(ps probe.Server, r hexa.HealthReporter) error {
	for _, d := range registry.Registry().Descriptors() {
		if health, ok := d.Instance.(hexa.Health); ok {
			r.AddToChecks(health)
		}
	}

	return tracer.Trace(ps.Run())
}
