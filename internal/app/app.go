package app

import (
	"github.com/kamva/tracer"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/model"
)

// appCore is implementation of the App
type appCore struct {
	tx App

	cfg   *config.Config
	store model.Store
	sp    base.ServiceProvider
}

// New returns new instance of the App
func New(sp base.ServiceProvider, store model.Store) (App, error) {
	return &appCore{
		cfg:   sp.Config().(*config.Config),
		store: store,
		sp:    sp,
	}, nil
}

// NewWithAllLayers returns a new app with all needed layers like DB transaction layer,...
func NewWithAllLayers(sp base.ServiceProvider, store model.Store) (App, error) {
	a, err := New(sp, store)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	// wrap app in SQL transaction layer. if we have more than one DB driver, we should
	// check driver type and create new instance TX layer for that DB type.
	//txLayer := NewSQLTxLayer(store, a)
	//a.(*appCore).tx = txLayer
	return NewTracingLayer(sp.OpenTelemetry().TracerProvider(), a), nil
}

var _ App = &appCore{}
