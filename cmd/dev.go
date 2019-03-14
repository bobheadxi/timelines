package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/store"
)

func newDevCommand() *cobra.Command {
	d := &cobra.Command{
		Use:              "dev",
		Hidden:           os.Getenv("MODE") != "development",
		PersistentPreRun: func(*cobra.Command, []string) { godotenv.Load() },
	}
	d.AddCommand(newRedisCommand(), newPGCommand())
	return d
}

func newPGCommand() *cobra.Command {
	var (
		envFiles []string
		pg       = &cobra.Command{
			Use: "pg",
			PersistentPreRun: func(*cobra.Command, []string) {
				if len(envFiles) > 0 {
					godotenv.Load(envFiles...)
				} else {
					godotenv.Load()
				}
			},
		}
	)
	pg.Flags().StringArrayVar(&envFiles, "env", nil, "")

	var (
		seed = &cobra.Command{
			Use: "seed",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zap.NewDevelopment()
				if err != nil {
					return err
				}
				var l = logger.Sugar().Named("dev.redis.reset")
				c, err := db.New(l, "integration_test", dev.DatabaseOptions)
				if err != nil {
					return err
				}
				return c.Repos().NewRepository(context.Background(),
					dev.GetTestInstallationID(),
					"bobheadxi", "calories")
			},
		}
	)

	pg.AddCommand(seed)
	return pg
}

func newRedisCommand() *cobra.Command {
	var (
		redis = &cobra.Command{
			Use: "redis",
		}
	)

	var (
		reset = &cobra.Command{
			Use: "reset",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zap.NewDevelopment()
				if err != nil {
					return err
				}
				var l = logger.Sugar().Named("dev.redis.reset")
				c, err := store.NewClient(l, dev.StoreOptions)
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
			Use: "seed",
			RunE: func(cmd *cobra.Command, args []string) error {
				logger, err := zap.NewDevelopment()
				if err != nil {
					return err
				}
				var l = logger.Sugar().Named("dev.redis.seed")
				c, err := store.NewClient(l, dev.StoreOptions)
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
