package commands

import (
	"time"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/model"
	"space.org/space/internal/registry"
	"space.org/space/internal/store"
	cachestore "space.org/space/internal/store/cache_layer"
	"space.org/space/internal/store/sqlstore"
)

func boot() (app.App, base.ServiceProvider, error) {
	sp, err := base.NewServiceProvider()
	if err != nil {
		return nil, nil, tracer.Trace(err)
	}

	cfg := sp.Config().(*config.Config)

	var s model.Store
	s, err = sqlstore.New(sp, cfg.DB)
	if err != nil {
		sp.Logger().Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		time.Sleep(time.Second) // Wait to flush the log buffer
		return nil, nil, tracer.Trace(err)
	}

	if sp.CacheProvider() != nil { // Add the cache layer.
		s = cachestore.New(sp, s)
	}

	s = store.NewTracingLayerStore("sql", sp.OpenTelemetry().TracerProvider(), s)

	a, err := app.NewWithAllLayers(sp, s)
	if err != nil {
		sp.Logger().Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		time.Sleep(time.Second)
		return nil, nil, tracer.Trace(err)
	}

	// Register store and app:
	registry.Register(registry.StoreService, s)
	registry.Register(registry.AppService, a)
	// Set global DB store on the model package:
	model.SetDBStore(s)

	return a, sp, nil
}
