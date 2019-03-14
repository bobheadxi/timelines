package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/worker"
)

func newWorkerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		devmode bool
	)
	c := &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := log.NewLogger(devmode, logpath)
			if err != nil {
				return err
			}
			// TODO: replace with real options
			godotenv.Load()
			return worker.Run(
				l.Named("worker"),
				newStopper(),
				worker.RunOpts{
					Store:    dev.StoreOptions,
					Database: dev.DatabaseOptions,
				})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8090", "")
	flags.StringVar(&logpath, "logpath", "", "")
	flags.BoolVar(&devmode, "dev", false, "")
	return c
}
