package commands

import (
	"context"
	"strconv"

	"github.com/kamva/tracer"
	"github.com/spf13/cobra"
	"t3.org/t3/internal/registry"
)

var excludeServcies = []string{
	registry.ServiceNameChannelHomes,
	registry.ServiceNameDispatcher,
}

var seederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "Seed data into database",
}

var SeederSeedCmd = &cobra.Command{
	Use:     "seed [num]",
	Short:   "Seed Data into DB",
	Example: "seed 100",
	Args:    cobra.ExactArgs(1),
	RunE:    withCtx(seedCmdF, registry.BaseServices(excludeServcies...)...),
}

func init() {
	seederCmd.AddCommand(SeederSeedCmd)

	rootCmd.AddCommand(seederCmd)
}

func seedCmdF(ctx context.Context, o *cmdOpts, _ *cobra.Command, args []string) error {
	err := registry.ProvideByNames(o.Registry)
	if err != nil {
		return tracer.Trace(err)
	}

	count, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(o.App.Seed(ctx, int32(count)))
}
