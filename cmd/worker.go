package cmd

import (
	"github.com/bobheadxi/projector/log"
	"github.com/bobheadxi/projector/worker"
	"github.com/spf13/cobra"
)

func newWorkerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		dev     bool
	)
	c := &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := log.NewLogger(dev, logpath)
			if err != nil {
				return err
			}
			return worker.Run(l, worker.RunOpts{
				Port: port,
			})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8090", "")
	flags.StringVar(&logpath, "logpath", "", "")
	flags.BoolVar(&dev, "dev", false, "")
	return c
}
