package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "space",
	Short: "space",
	Long:  "space",
}

func Run() error {
	return rootCmd.Execute()
}
