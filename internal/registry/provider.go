package registry

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

type Provider func(r hexa.ServiceRegistry) error

func ProvideServices(r hexa.ServiceRegistry, providers map[string]Provider) error {
	// Append services that are not in the sort list to the services list that we should provide.
	servicesPriority := bootPriority()
	names := append(servicesPriority, gutil.Sub(servicesPriority, serviceNamesFromMap(providers))...)

	for _, name := range names {
		if p, ok := providers[name]; ok {
			if err := Provide(r, p); err != nil {
				return tracer.Trace(err)
			}
		}
	}

	return nil
}

func Provide(r hexa.ServiceRegistry, p Provider) error {
	return tracer.Trace(p(r))
}

func serviceNamesFromMap(providers map[string]Provider) []string {
	names := make([]string, len(providers))
	var i int
	for name := range providers {
		names[i] = name
		i++
	}
	return names
}
