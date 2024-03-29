package commands

import (
	"context"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"t3.org/t3/infra"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
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

func withApp(cmdF WithAppHandler, serviceNames ...string) func(cmd *cobra.Command, args []string) error {
	// serviceNames never can not be empty (because at least the app itself should be in the list),
	// so we'll check if it's empty, we'll set our default services list.
	if len(serviceNames) == 0 {
		serviceNames = registry.BaseServices()
	}
	return func(cmd *cobra.Command, args []string) error {
		r := sr.New()
		err := registry.ProvideByNames(r, serviceNames...)
		if err != nil {
			return tracer.Trace(err)
		}

		if err := r.Boot(); err != nil {
			return tracer.Trace(err)
		}

		timeout := time.Second * 30
		go sr.ShutdownBySignals(r, timeout)      //nolint
		defer sr.ShutdownWithTimeout(r, timeout) //nolint

		return cmdF(&cmdOpts{
			Registry: r,
			Cfg:      r.Service(registry.ServiceNameConfig).(*config.Config),
			App:      r.Service(registry.ServiceNameApp).(app.App),
		}, cmd, args)
	}
}

//nolint:unused
func withCtx(cmdF WithCtxHandler, serviceNames ...string) func(cmd *cobra.Command, args []string) error {
	return withApp(func(o *cmdOpts, cmd *cobra.Command, args []string) error {
		s := services.New(o.Registry)
		u := infra.NewServiceUser(infra.UserIdCommandLine)
		ctx := hexa.NewContext(context.Background(), hexa.ContextParams{
			CorrelationId:  gutil.UUID(),
			Locale:         "en-US",
			User:           u,
			BaseLogger:     s.Logger(),
			BaseTranslator: s.Translator(),
		})
		return cmdF(ctx, o, cmd, args)
	}, serviceNames...)
}
