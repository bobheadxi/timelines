package cmd

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/bobheadxi/projector/dev"
	"github.com/bobheadxi/projector/store"
)

func newDevCommand() *cobra.Command {
	d := &cobra.Command{
		Use:              "dev",
		Hidden:           true,
		PersistentPreRun: func(*cobra.Command, []string) { godotenv.Load() },
	}
	d.AddCommand(newSeedCommand())
	return d
}

func newSeedCommand() *cobra.Command {
	seed := &cobra.Command{
		Use: "seed",
	}
	seedRedis := &cobra.Command{
		Use: "redis",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := zap.NewDevelopment()
			if err != nil {
				return err
			}
			var l = logger.Sugar().Named("dev.seed.redis")
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

	seed.AddCommand(seedRedis)
	return seed
}
