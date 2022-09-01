package commands

import (
	"context"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	hexarobfig "github.com/kamva/hexa-job/robfig"
	"github.com/kamva/tracer"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/router/crons"
)

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "Manage cron jobs",
}

var cronRunCommand = &cobra.Command{
	Use:     "run",
	Short:   "Run cron jobs",
	Example: "cron run",
	RunE:    withApp(runCron),
}

func init() {
	rootCmd.AddCommand(cronCommand)
	cronCommand.AddCommand(cronRunCommand)
}

func bootCron(cfg *config.Config, sp base.ServiceProvider, a app.App) (hjob.CronJobs, error) {
	u, err := hexa.NewGuest().SetMeta(hexa.UserMetaKeyName, "cron_job")
	if err != nil {
		return nil, tracer.Trace(err)
	}
	ctxGen := func(c context.Context) context.Context {
		return hexa.NewContext(c, hexa.ContextParams{
			Locale:         "en",
			User:           u,
			CorrelationId:  gutil.UUID(),
			BaseLogger:     sp.Logger(),
			BaseTranslator: sp.Translator(),
		})
	}

	cronJobs := hexarobfig.New(ctxGen, cron.New())

	if err := crons.RegisterCronJobs(cronJobs, sp, a); err != nil {
		return nil, tracer.Trace(err)
	}

	registry.Register(registry.CronService, cronJobs)
	return cronJobs, nil
}

func runCron(o *cmdOpts, cmd *cobra.Command, args []string) error {
	a, sp, cfg := o.App, o.SP, o.Cfg

	cronJobs, err := bootCron(cfg, sp, a)
	if err != nil {
		return tracer.Trace(err)
	}

	// Run healthChecker server:
	if err := runProbeServer(o.SP); err != nil {
		return tracer.Trace(err)
	}

	app.Banner("Space cron jobs")
	return tracer.Trace(cronJobs.Run())
}
