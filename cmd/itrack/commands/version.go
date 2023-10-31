package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"space.org/space/internal/app"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Shows the app version",
	Example: "version",
	//Args:    cobra.ExactArgs(1),
	RunE: printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(_ *cobra.Command, _ []string) error {
	fmt.Print(app.Version)
	return nil
}
