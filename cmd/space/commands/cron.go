package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
	"space.org/space/internal/registry/services"
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

func runCron(o *cmdOpts, cmd *cobra.Command, args []string) error {
	s := services.New(o.Registry)
	if err := registry.Provide(o.Registry, provider.CronProvider); err != nil {
		return tracer.Trace(err)
	}

	// Run healthChecker server:
	if err := runProbeServer(o.Registry, s.ProbeServer(), s.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	app.Banner("Space cron jobs")
	return tracer.Trace(s.CronJobs().Run())
}
