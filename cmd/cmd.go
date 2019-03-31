package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// Initialize bootstraps the giveen command with necessary child commands
func Initialize(cmd *cobra.Command) {
	cmd.AddCommand(
		newServerCmd(),
		newWorkerCmd(),
		newDevCommand())
}

func newStopper() chan bool {
	var (
		stopper = make(chan bool)
		signals = make(chan os.Signal)
	)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signals
		stopper <- true
	}()

	return stopper
}
