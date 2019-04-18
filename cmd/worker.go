package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd/monitoring"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/worker"
)

func newWorkerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		workers int
		profile bool
		devMode bool
	)
	c := &cobra.Command{
		Use:   "worker",
		Short: "spin up a Timelines worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			meta := config.NewBuildMeta()
			l, err := log.NewLogger(devMode, logpath)
			if err != nil {
				return err
			}
			l = l.With("build.version", meta.Commit)

			if profile {
				l.Info("setting up profiling")
				if err := monitoring.StartProfiler(l, "timelines-worker", meta); err != nil {
					return fmt.Errorf("failed to start profiler: %v", err)
				}
			}

			storeCfg := config.NewStoreConfig()
			dbCfg := config.NewDatabaseConfig()
			if devMode {
				storeCfg = dev.StoreOptions
				dbCfg = dev.DatabaseOptions
			}

			return worker.Run(
				l.Named("worker"),
				newStopper(),
				worker.RunOpts{
					Workers:  workers,
					Store:    storeCfg,
					Database: dbCfg,
				})
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8090", "port to serve worker API on")
	flags.StringVar(&logpath, "logpath", "", "path to log dump")
	flags.IntVar(&workers, "workers", 3, "number of workers to spin up")
	flags.BoolVar(&profile, "profile", false, "enable profiling")
	flags.BoolVar(&devMode, "dev", false, "toggle dev mode")
	return c
}
