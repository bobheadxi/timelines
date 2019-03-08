package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

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
