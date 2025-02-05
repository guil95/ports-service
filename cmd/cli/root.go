package cli

import (
	"github.com/spf13/cobra"
	"log/slog"
)

var RootCmd = &cobra.Command{
	Use:          "ports-service",
	Short:        "Ports Service is a service to manage ports data",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("No subcommand provided, running 'server' by default")
		return ServeCmd.RunE(cmd, args)
	},
}
