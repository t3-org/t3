package registry

import (
	"fmt"
	"sync"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type Provider func(r hexa.ServiceRegistry) error

var providerToServiceMap = make(map[string]string) // map[providerName]serviceName.
var providers = make(map[string]Provider)
var plistLock sync.Mutex

func AddProvider(providerName string, serviceName string, p Provider) {
	plistLock.Lock()
	defer plistLock.Unlock()
	if _, ok := providerToServiceMap[providerName]; ok {
		hlog.Warn("the provider is already provided, overwrite it", hlog.String("provider", providerName))
	}

	providerToServiceMap[providerName] = serviceName
	providers[providerName] = p
}

// Providers returns provides.
// The return param is map[serviceName]Provider.
// names is list of providers name.
func Providers(names ...string) (map[string]Provider, error) {
	plistLock.Lock()
	defer plistLock.Unlock()
	m := make(map[string]Provider)
	for _, name := range names {
		svcName, ok := providerToServiceMap[name]
		if !ok {
			return nil, fmt.Errorf("can not find provider with name: %s", name)
		}
		m[svcName] = providers[name]
	}
	return m, nil
}

func ProviderByName(name string) Provider {
	plistLock.Lock()
	defer plistLock.Unlock()
	return providers[name]
}

// ProvideByProviders provides services using providers. providers param is map[serviceName]Provider.
func ProvideByProviders(r hexa.ServiceRegistry, providers map[string]Provider) error {
	// Append services that are not in the sort list to the services list that we should provide.
	servicesPriority := bootPriority()
	names := append(servicesPriority, gutil.Sub(servicesPriority, mapKeys(providers))...)

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

// ProvideByName provides a service by the provider name.
func ProvideByName(r hexa.ServiceRegistry, name string) error {
	return Provide(r, ProviderByName(name))
}

// ProvideByNames provides services using their corresponding provider name.
func ProvideByNames(r hexa.ServiceRegistry, names ...string) error {
	providers, err := Providers(names...)
	if err != nil {
		return tracer.Trace(err)
	}
	return ProvideByProviders(r, providers)
}

func mapKeys(providers map[string]Provider) []string {
	names := make([]string, len(providers))
	var i int
	for name := range providers {
		names[i] = name
		i++
	}
	return names
}
