package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/cmd/monitoring"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/server"
	"github.com/bobheadxi/zapx"
	"github.com/bobheadxi/zapx/zgcp"
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
		Run: func(cmd *cobra.Command, args []string) {
			meta := config.NewBuildMeta()
			logger, err := zapx.New(logpath, devMode)
			if err != nil {
				println("logger: " + err.Error())
				os.Exit(1)
			}
			l := logger.Sugar().Named(monitor.Service).With("build.version", meta.AnnotatedCommit(devMode))
			l.Infof("preparing to start %s", monitor.Service)

			if monitor.Profile {
				if err := monitoring.StartProfiler(l, monitor.Service, meta, devMode); err != nil {
					l.Fatalf("failed to start profiler: %v", err)
				}
			}
			if monitor.Errors {
				cloud := config.NewCloudConfig()
				errorLogger, err := zgcp.NewErrorReportingLogger(
					l.Desugar(),
					zgcp.ServiceConfig{
						ProjectID: cloud.GCP.ProjectID,
						Name:      monitor.Service,
						Version:   meta.AnnotatedCommit(devMode),
					},
					zgcp.Fields{
						UserKey: config.LogKeyRID,
					},
					devMode,
					config.NewGCPConnectionOptions()...)
				if err != nil {
					l.Fatalf("failed to attach error logger: %v", err)
				}
				l = errorLogger.Sugar()
			}

			storeCfg := config.NewStoreConfig()
			dbCfg := config.NewDatabaseConfig()
			if devMode {
				storeCfg = dev.StoreOptions
				dbCfg = dev.DatabaseOptions
			}

			defer l.Sync()
			if err := server.Run(
				l,
				newStopper(),
				server.RunOpts{
					Port:     port,
					Store:    storeCfg,
					Database: dbCfg,
					Meta:     meta,
					Dev:      devMode,
				},
			); err != nil {
				l.Fatalf("server exited with error: %v", err)
			}
		},
	}
	flags := c.Flags()
	flags.StringVarP(&port, "port", "p", "8080", "port to serve API on")
	flags.StringVar(&logpath, "logpath", "", "path to log dump")
	monitor.Attach(flags, "timelines-server")
	flags.BoolVar(&devMode, "dev", false, "toggle dev mode")
	return c
}
