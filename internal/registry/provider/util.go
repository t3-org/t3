package provider

import (
	"github.com/kamva/hexa"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
)

func conf(r hexa.ServiceRegistry) *config.Config {
	return r.Service(registry.ServiceNameConfig).(*config.Config)
}
