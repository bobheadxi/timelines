package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bobheadxi/projector/log"
	"github.com/bobheadxi/projector/server"
)

func newServerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		dev     bool
	)
	c := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := log.NewLogger(dev, logpath)
			if err != nil {
				return err
			}
			r := server.NewResolver()
			return server.Run(l, r, server.RunOpts{
				Port: port,
			})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8080", "")
	flags.StringVar(&logpath, "logpath", "", "")
	flags.BoolVar(&dev, "dev", false, "")
	return c
}
