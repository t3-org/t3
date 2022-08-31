package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "shield",
	Short: "Shield IAM",
	Long:  "Shield IAM",
}

func Run() error {
	return rootCmd.Execute()
}
