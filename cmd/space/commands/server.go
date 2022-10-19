package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/services"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Management of the http server",
}

var serverListenCMD = &cobra.Command{
	Use:     "listen",
	Short:   "Run the http server",
	Example: "listen",
	RunE:    withApp(serverCmdF),
}

func init() {
	serverCmd.AddCommand(serverListenCMD)

	rootCmd.AddCommand(serverCmd)
}

func serverCmdF(o *cmdOpts, _ *cobra.Command, _ []string) error {
	s := services.New(o.Registry)
	if err := registry.ProvideByName(o.Registry, registry.ServiceNameHttpServer); err != nil {
		return tracer.Trace(err)
	}

	if err := runProbeServer(o.Registry, s.ProbeServer(), s.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	// Start server
	app.Banner("Space")
	done, err := s.HttpServer().Run()
	if err != nil {
		return tracer.Trace(err)
	}
	return tracer.Trace(<-done)
}
