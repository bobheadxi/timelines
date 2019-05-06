package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/bobheadxi/timelines/analysis/testdata"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/store"
	"go.bobheadxi.dev/zapx"
)

func newDevCommand() *cobra.Command {
	d := &cobra.Command{
		Use:    "dev",
		Short:  "handy utility commands for development use",
		Hidden: os.Getenv("MODE") != "development",
	}
	d.AddCommand(newRedisCommand(), newPGCommand())
	return d
}

func newPGCommand() *cobra.Command {
	var (
		fromEnv bool
		pg      = &cobra.Command{
			Use:   "pg",
			Short: "postgres database utilities",
		}
	)
	pg.PersistentFlags().BoolVar(&fromEnv, "from-env", false, "use env settings instead of hardcoded dev settings")

	var (
		seed = &cobra.Command{
			Use:   "seed",
			Short: "seed postgres database with test data",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zapx.New("", true)
				if err != nil {
					return err
				}

				var (
					l = logger.Sugar().Named("dev.pg.seed")
					c *db.Database
				)
				if fromEnv {
					l.Info("loading config from env")
					c, err = db.New(l, "integration_test", config.NewDatabaseConfig())
				} else {
					l.Info("loading local development config")
					c, err = db.New(l, "integration_test", dev.DatabaseOptions)
				}
				if err != nil {
					return err
				}

				c.Repos().NewRepository(
					context.Background(),
					dev.GetTestInstallationID(),
					dev.SampleRepo)
				repo, err := c.Repos().GetRepository(
					context.Background(),
					dev.SampleRepo.Host,
					dev.SampleRepo.Owner, dev.SampleRepo.Name)
				if err != nil {
					return err
				}
				c.Repos().InsertGitBurndownResult(context.Background(),
					repo.ID,
					testdata.Meta,
					testdata.Burndown)
				return nil
			},
		}
	)

	pg.AddCommand(seed)
	return pg
}

func newRedisCommand() *cobra.Command {
	var (
		fromEnv bool
		redis   = &cobra.Command{
			Use:   "redis",
			Short: "redis store utilities",
		}
	)
	redis.PersistentFlags().BoolVar(&fromEnv, "from-env", false, "use env settings instead of hardcoded dev settings")

	var (
		reset = &cobra.Command{
			Use:   "reset",
			Short: "drop everything in store",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zapx.New("", true)
				if err != nil {
					return err
				}

				var (
					l = logger.Sugar().Named("dev.redis.reset")
					c *store.Client
				)
				if fromEnv {
					l.Info("loading config from env")
					c, err = store.NewClient(l, "dev.redis.reset", config.NewStoreConfig())
				} else {
					l.Info("loading local development config")
					c, err = store.NewClient(l, "dev.redis.reset", dev.StoreOptions)
				}
				if err != nil {
					return err
				}

				c.Reset()
				return nil
			},
		}
	)

	var (
		seed = &cobra.Command{
			Use:   "seed",
			Short: "seed store with test data",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zapx.New("", true)
				if err != nil {
					return err
				}

				var (
					l = logger.Sugar().Named("dev.redis.seed")
					c *store.Client
				)
				if fromEnv {
					l.Info("loading config from env")
					c, err = store.NewClient(l, "dev.redis.seed", config.NewStoreConfig())
				} else {
					l.Info("loading local development config")
					c, err = store.NewClient(l, "dev.redis.seed", dev.StoreOptions)
				}
				if err != nil {
					return err
				}

				// queue job and add state entry
				id, _ := uuid.NewUUID()
				if err = c.RepoJobs().Queue(&store.RepoJob{
					ID:             id,
					Owner:          "bobheadxi",
					Repo:           "calories",
					InstallationID: os.Getenv("GITHUB_TEST_INSTALLTION"),
				}); err != nil {
					return fmt.Errorf("failed to queue job: %s", err.Error())
				}

				return nil
			},
		}
	)

	redis.AddCommand(seed, reset)
	return redis
}
