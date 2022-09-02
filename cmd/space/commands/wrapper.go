package commands

import (
	"context"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/infra"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
)

type cmdOpts struct {
	Cfg *config.Config
	SP  base.ServiceProvider
	App app.App
}

// WithAppHandler gets the app, service-provider and config as params to handle the command
type WithAppHandler func(o *cmdOpts, cmd *cobra.Command, args []string) error

// WithCtxHandler gets the hexa context, app, service-provider and config as params to handle the command
type WithCtxHandler func(ctx context.Context, o *cmdOpts, cmd *cobra.Command, args []string) error

func withApp(cmdF WithAppHandler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := registry.ProvideServices(registry.Registry(), provider.Providers(registry.BaseServices()))
		if err != nil {
			return tracer.Trace(err)
		}

		sp := registry.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
		a := registry.Service(registry.ServiceNameApp).(app.App)

		cfg := sp.Config().(*config.Config)

		timeout := time.Second * 30
		go sr.ShutdownBySignals(registry.Registry(), timeout)
		defer registry.Shutdown(timeout)

		return cmdF(&cmdOpts{Cfg: cfg, SP: sp, App: a}, cmd, args)
	}
}

func withCtx(cmdF WithCtxHandler) func(cmd *cobra.Command, args []string) error {
	return withApp(func(o *cmdOpts, cmd *cobra.Command, args []string) error {
		u := infra.NewServiceUser(infra.UserIdCommandLine)
		ctx := hexa.NewContext(nil, hexa.ContextParams{
			CorrelationId:  gutil.UUID(),
			Locale:         "en-US",
			User:           u,
			BaseLogger:     o.SP.Logger(),
			BaseTranslator: o.SP.Translator(),
		})
		return cmdF(ctx, o, cmd, args)
	})
}
