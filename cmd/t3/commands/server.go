package commands

import (
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
	"t3.org/t3/internal/router/matrix"
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
	err := registry.ProvideByNames(
		o.Registry,
		registry.ServiceNameHttpServer,
		registry.ServiceNameMatrixBotServer,
	)

	if err != nil {
		return tracer.Trace(err)
	}

	if err := runProbeServer(o.Registry, s.ProbeServer(), s.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	// Start server
	app.Banner("Space")

	mxDone, err := o.Registry.Service(registry.ServiceNameMatrixBotServer).(*matrix.Server).Run()
	if err != nil {
		return tracer.Trace(err)
	}

	done, err := s.HttpServer().Run()
	if err != nil {
		return tracer.Trace(err)
	}

	return gutil.AnyErr(tracer.Trace(<-mxDone), tracer.Trace(<-done))
}
