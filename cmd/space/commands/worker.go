package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
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

func runWorker(o *cmdOpts, cmd *cobra.Command, args []string) error {
	services := base.NewServices(o.Registry)
	if err := registry.Provide(o.Registry, provider.WorkerProvider); err != nil {
		return tracer.Trace(err)
	}

	// Run healthChecker server:
	if err := runProbeServer(o.Registry, services.ProbeServer(), services.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	app.Banner("Space worker")
	return tracer.Trace(services.Worker().Run())
}
