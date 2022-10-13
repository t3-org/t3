package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/services"
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

func runWorker(o *cmdOpts, _ *cobra.Command, _ []string) error {
	s := services.New(o.Registry)
	if err := registry.ProvideByName(o.Registry, registry.ServiceNameWorker); err != nil {
		return tracer.Trace(err)
	}

	// Run healthChecker server:
	if err := runProbeServer(o.Registry, s.ProbeServer(), s.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	app.Banner("Space worker")
	return tracer.Trace(s.Worker().Run())
}
