package registry

import (
	"context"
	"time"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
)

var r = sr.New()

func Registry() hexa.ServiceRegistry {
	return r
}

func Register(name string, instance interface{}) {
	r.Register(name, instance)
}

func Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return tracer.Trace(r.Shutdown(ctx))
}

func Service(name string) any {
	return r.Service(name)
}
