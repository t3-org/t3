package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "t3",
	Short: "t3",
	Long:  "t3",
}

func Run() error {
	return rootCmd.Execute()
}
