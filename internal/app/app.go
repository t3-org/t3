package app

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/tracer"
	"t3.org/t3/internal/service/channel"

	"t3.org/t3/internal/config"
	"t3.org/t3/internal/model"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
)

// appCore is implementation of the App
type appCore struct {
	cfg        *config.Config
	store      model.Store
	dispatcher *channel.Dispatcher
}

// New returns new instance of the App
func New(r hexa.ServiceRegistry, store model.Store) (App, error) {
	s := services.New(r)

	return &appCore{
		cfg:        s.Config(),
		store:      store,
		dispatcher: s.Dispatcher(),
	}, nil
}

// NewWithAllLayers returns a new app with all needed layers like DB transaction layer,...
func NewWithAllLayers(r hexa.ServiceRegistry, store model.Store) (App, error) {
	a, err := New(r, store)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	// wrap app in SQL transaction layer. if we have more than one DB driver, we should
	// check driver type and create new instance TX layer for that DB type.
	//txLayer := NewSQLTxLayer(store, a)
	//a.(*appCore).tx = txLayer
	return NewTracingLayer(r.Service(registry.ServiceNameOpenTelemetry).(htel.OpenTelemetry).TracerProvider(), a), nil
}

var _ App = &appCore{}
