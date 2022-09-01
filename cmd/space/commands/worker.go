package commands

import (
	"github.com/hibiken/asynq"
	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/hexa-job/hsynq"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/router/jobs"
)

var workerCommand = &cobra.Command{
	Use:   "worker",
	Short: "Manage worker",
}

var workerRunCommand = &cobra.Command{
	Use:     "run",
	Short:   "Run to do tasks",
	Example: "worker run",
	//Args:    cobra.ExactArgs(1),
	RunE: withApp(runWorker),
}

func init() {
	rootCmd.AddCommand(workerCommand)
	workerCommand.AddCommand(workerRunCommand)
}

func bootWorker(cfg *config.Config, sp base.ServiceProvider, a app.App) (hjob.Worker, error) {
	srv := asynq.NewServer(
		cfg.AsynqConfig.RedisOpts(),
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: cfg.AsynqConfig.WorkerConcurrency,
			// Optionally specify multiple queues with different priority.
			Queues: cfg.AsynqConfig.Queues(),
		},
	)

	w := hsynq.NewWorker(srv, hexa.NewContextPropagator(sp.Logger(), sp.Translator()), hsynq.NewJsonTransformer())
	if err := jobs.RegisterJobs(w, sp, a); err != nil {
		return nil, tracer.Trace(err)
	}
	registry.Register(registry.WorkerService, w)
	return w, nil
}

func runWorker(o *cmdOpts, cmd *cobra.Command, args []string) error {
	a, sp, cfg := o.App, o.SP, o.Cfg

	worker, err := bootWorker(cfg, sp, a)
	if err != nil {
		return tracer.Trace(err)
	}

	// Run healthChecker server:
	if err := runProbeServer(o.SP); err != nil {
		return tracer.Trace(err)
	}

	app.Banner("Space worker")
	return tracer.Trace(worker.Run())
}
