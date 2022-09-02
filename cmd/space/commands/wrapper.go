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
	Registry hexa.ServiceRegistry
	Cfg      *config.Config
	App      app.App
}

// WithAppHandler gets the app, service-provider and config as params to handle the command
type WithAppHandler func(o *cmdOpts, cmd *cobra.Command, args []string) error

// WithCtxHandler gets the hexa context, app, service-provider and config as params to handle the command
type WithCtxHandler func(ctx context.Context, o *cmdOpts, cmd *cobra.Command, args []string) error

func withApp(cmdF WithAppHandler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		r := sr.New()
		err := registry.ProvideServices(r, provider.Providers(registry.BaseServices()))
		if err != nil {
			return tracer.Trace(err)
		}

		timeout := time.Second * 30
		go sr.ShutdownBySignals(r, timeout)
		defer registry.Shutdown(r, timeout)

		return cmdF(&cmdOpts{
			Registry: r,
			Cfg:      r.Service(registry.ServiceNameConfig).(*config.Config),
			App:      r.Service(registry.ServiceNameApp).(app.App),
		}, cmd, args)
	}
}

func withCtx(cmdF WithCtxHandler) func(cmd *cobra.Command, args []string) error {
	return withApp(func(o *cmdOpts, cmd *cobra.Command, args []string) error {
		services := base.NewServices(o.Registry)
		u := infra.NewServiceUser(infra.UserIdCommandLine)
		ctx := hexa.NewContext(nil, hexa.ContextParams{
			CorrelationId:  gutil.UUID(),
			Locale:         "en-US",
			User:           u,
			BaseLogger:     services.Logger(),
			BaseTranslator: services.Translator(),
		})
		return cmdF(ctx, o, cmd, args)
	})
}
