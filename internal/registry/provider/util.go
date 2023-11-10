package provider

import (
	"github.com/kamva/hexa"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/registry"
)

func conf(r hexa.ServiceRegistry) *config.Config {
	return r.Service(registry.ServiceNameConfig).(*config.Config)
}
