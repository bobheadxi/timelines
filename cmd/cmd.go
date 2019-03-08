package cmd

import (
	"github.com/spf13/cobra"
)

// Initialize bootstraps the giveen command with necessary child commands
func Initialize(cmd *cobra.Command) {
	cmd.AddCommand(
		newServerCmd(),
		newWorkerCmd())
}
