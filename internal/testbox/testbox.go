package testbox

import (
	"time"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
	"space.org/space/internal/registry"
)

var testbox *TestBox

func SetGlobal(tbox *TestBox) {
	testbox = tbox
}

func Global() *TestBox {
	return testbox
}

var shutdownTimeout = time.Second * 30

// TestBox provides a configurable test environment.
type TestBox struct {
	r         hexa.ServiceRegistry
	providers map[string]registry.Provider
}

func New(providers map[string]registry.Provider) *TestBox {
	return &TestBox{
		r:         sr.New(),
		providers: providers,
	}
}

func (t *TestBox) Setup() error {
	if err := registry.ProvideServices(t.r, t.providers); err != nil {
		return tracer.Trace(err)
	}

	if err := t.r.Boot(); err != nil {
		return tracer.Trace(err)
	}

	go sr.ShutdownBySignals(t.r, shutdownTimeout) //nolint
	return nil
}

func (t *TestBox) Teardown() error {
	return registry.Shutdown(t.r, shutdownTimeout)
}

func (t *TestBox) TeardownIfPanic() {
	err := recover()
	if err == nil {
		return
	}

	hlog.Error("recovered from test")
	if teardownErr := t.Teardown(); teardownErr != nil {
		hlog.Error("can not teardown testbox completely", hlog.ErrStack(teardownErr))
	}
	panic(err)
}

func (t *TestBox) Registry() hexa.ServiceRegistry {
	return t.r
}
