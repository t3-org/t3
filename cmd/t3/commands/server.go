package commands

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
	"t3.org/t3/internal/service/channel"
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
	)

	if err != nil {
		return tracer.Trace(err)
	}

	if err := runProbeServer(o.Registry, s.ProbeServer(), s.HealthReporter()); err != nil {
		return tracer.Trace(err)
	}

	// Start server
	app.Banner("T3")

	chHomes := o.Registry.Service(registry.ServiceNameChannelHomes).(map[string]channel.Home)
	channels := make(map[string]<-chan error)
	for name, home := range chHomes {
		if runnable := home.(hexa.Runnable); runnable != nil {
			hlog.Info("running channel home", hlog.String("name", name))
			closeCh, err := runnable.Run()
			if err != nil {
				return tracer.Trace(err)
			}
			channels[name] = closeCh
		}
	}

	done, err := s.HttpServer().Run()
	if err != nil {
		return tracer.Trace(err)
	}

	if err := <-done; err != nil {
		hlog.Error("error on server", hlog.Err(err))
	}

	// Waiting till close all runnable channel homes:
	for name, ch := range channels {
		if err := <-ch; err != nil {
			hlog.Error("error on channel", hlog.String("name", name), hlog.Err(err))
		}
	}
	return nil
}
