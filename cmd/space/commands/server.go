package commands

import (
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
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

func serverCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	sp := base.NewServiceProvider(o.Registry)
	if err := registry.Provide(o.Registry, provider.HttpServerProvider); err != nil {
		return tracer.Trace(err)
	}

	if err := runProbeServer(o.Registry, sp.ProbeServer(), sp.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	// Start server
	app.Banner("Space")
	return tracer.Trace(sp.HttpServer().Run())
}
