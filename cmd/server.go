package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/server"
)

func newServerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		devMode bool
	)
	c := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := log.NewLogger(devMode, logpath)
			if err != nil {
				return err
			}

			storeCfg := config.NewStoreConfig()
			dbCfg := config.NewDatabaseConfig()
			if devMode {
				storeCfg = dev.StoreOptions
				dbCfg = dev.DatabaseOptions
			}

			return server.Run(
				l.Named("server"),
				newStopper(),
				server.RunOpts{
					Port:     port,
					Store:    storeCfg,
					Database: dbCfg,
				})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8080", "")
	flags.StringVar(&logpath, "logpath", "", "")
	flags.BoolVar(&devMode, "dev", false, "")
	return c
}
