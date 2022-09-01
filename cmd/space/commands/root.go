package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "Space",
	Short: "Space",
	Long:  "Space",
}

func Run() error {
	return rootCmd.Execute()
}
