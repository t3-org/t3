package registry

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"space.org/space/internal/config"
)

type Provider func(r hexa.ServiceRegistry, cfg *config.Config) error

func ProvideServices(r hexa.ServiceRegistry, providers map[string]Provider) error {
	// Append services that are not in the sort list to the services list that we should provide.
	servicesPriority := bootPriority()
	names := append(servicesPriority, gutil.Sub(servicesPriority, serviceNamesFromMap(providers))...)

	// Initialize configs:
	cfg, err := config.New()
	if err != nil {
		return tracer.Trace(err)
	}
	config.SetDefaultConfig(cfg)
	r.Register(ServiceNameConfig, cfg)

	for _, name := range names {
		if p, ok := providers[name]; ok {
			if err := p(r, cfg); err != nil {
				return tracer.Trace(err)
			}
		}
	}

	return nil
}

func Provide(r hexa.ServiceRegistry, p Provider) error {
	cfg := r.Service(ServiceNameConfig).(*config.Config)
	return tracer.Trace(p(r, cfg))
}

func serviceNamesFromMap(providers map[string]Provider) []string {
	names := make([]string, len(providers))
	var i int
	for name, _ := range providers {
		names[i] = name
		i++
	}
	return names
}
