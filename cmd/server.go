package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/monitoring"
	"github.com/bobheadxi/timelines/server"
)

func newServerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		monitor = &monitoring.Flags{}
		devMode bool
	)
	c := &cobra.Command{
		Use:   "server",
		Short: "spin up the core Timelines server",
		RunE: func(cmd *cobra.Command, args []string) error {
			meta := config.NewBuildMeta()
			l, err := log.NewLogger(devMode, logpath)
			if err != nil {
				return err
			}
			l = l.Named(monitor.Service).With("build.version", meta.AnnotatedCommit(devMode))

			if monitor.Profile {
				if err := monitoring.StartProfiler(l, monitor.Service, meta, devMode); err != nil {
					return fmt.Errorf("failed to start profiler: %v", err)
				}
			}
			if monitor.Errors {
				if l, err = monitoring.AttachErrorLogging(l.Desugar(), monitor.Service, meta, devMode); err != nil {
					return fmt.Errorf("failed to attach error logger: %v", err)
				}
			}

			storeCfg := config.NewStoreConfig()
			dbCfg := config.NewDatabaseConfig()
			if devMode {
				storeCfg = dev.StoreOptions
				dbCfg = dev.DatabaseOptions
			}

			defer l.Sync()
			return server.Run(
				l,
				newStopper(),
				server.RunOpts{
					Port:     port,
					Store:    storeCfg,
					Database: dbCfg,
					Build:    meta.Commit,
				})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8080", "port to serve API on")
	flags.StringVar(&logpath, "logpath", "", "path to log dump")
	monitor.Attach(flags, "timelines-server")
	flags.BoolVar(&devMode, "dev", false, "toggle dev mode")
	return c
}
