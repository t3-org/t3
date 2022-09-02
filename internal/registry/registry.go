package registry

import (
	"context"
	"time"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

func Shutdown(r hexa.ServiceRegistry, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return tracer.Trace(r.Shutdown(ctx))
}
