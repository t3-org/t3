package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
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
	done, err := s.Worker().Run()
	if err != nil {
		return tracer.Trace(err)
	}
	return tracer.Trace(<-done)
}
