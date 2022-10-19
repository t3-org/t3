package commands

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/probe"
	"github.com/kamva/tracer"
)

func runProbeServer(r hexa.ServiceRegistry, ps probe.Server, reporter hexa.HealthReporter) error {
	for _, d := range r.Descriptors() {
		if health, ok := d.Instance.(hexa.Health); ok {
			reporter.AddToChecks(health)
		}
	}

	_, err := ps.Run()
	return tracer.Trace(err)
}
