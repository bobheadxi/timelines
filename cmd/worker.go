package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/monitoring"
	"github.com/bobheadxi/timelines/worker"
)

func newWorkerCmd() *cobra.Command {
	var (
		port    string
		logpath string
		workers int
		monitor = &monitoring.Flags{}
		devMode bool
	)
	c := &cobra.Command{
		Use:   "worker",
		Short: "spin up a Timelines worker",
		Run: func(cmd *cobra.Command, args []string) {
			meta := config.NewBuildMeta()
			l, err := log.NewLogger(devMode, logpath)
			if err != nil {
				println("logger: " + err.Error())
				os.Exit(1)
			}
			l = l.Named(monitor.Service).With("build.version", meta.AnnotatedCommit(devMode))
			l.Infof("preparing to start %s", monitor.Service)

			if monitor.Profile {
				if err := monitoring.StartProfiler(l, monitor.Service, meta, devMode); err != nil {
					l.Fatalf("failed to start profiler: %v", err)
				}
			}
			if monitor.Errors {
				errorLogger, err := monitoring.AttachErrorLogging(l.Desugar(), monitor.Service, meta, devMode)
				if err != nil {
					l.Fatalf("failed to attach error logger: %v", err)
				}
				l = errorLogger
			}

			storeCfg := config.NewStoreConfig()
			dbCfg := config.NewDatabaseConfig()
			if devMode {
				storeCfg = dev.StoreOptions
				dbCfg = dev.DatabaseOptions
			}

			defer l.Sync()
			if err := worker.Run(
				l,
				newStopper(),
				worker.RunOpts{
					Workers:  workers,
					Store:    storeCfg,
					Database: dbCfg,
				},
			); err != nil {
				l.Fatalf("worker exited with error: %v", err)
			}
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8090", "port to serve worker API on")
	flags.StringVar(&logpath, "logpath", "", "path to log dump")
	flags.IntVar(&workers, "workers", 3, "number of workers to spin up")
	monitor.Attach(flags, "timelines-worker")
	flags.BoolVar(&devMode, "dev", false, "toggle dev mode")
	return c
}
