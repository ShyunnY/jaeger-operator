package cmd

import "github.com/spf13/cobra"

// GetRootCommand returns the root cobra command to be executed by main.
func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "jaeger-operator",
		Short: "Jaeger Operator",
		Long:  "The Kubernetes operator for Jaeger",
	}

	rootCmd.AddCommand(GetServerCommand())

	return rootCmd
}
