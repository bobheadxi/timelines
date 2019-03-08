package cmd

import (
	"errors"

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
			return errors.New("not implemented")
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8080", "")
	flags.StringVar(&logpath, "logpath", "", "")
	flags.BoolVar(&dev, "dev", false, "")
	return c
}
