package db_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/timelines/analysis/testdata"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/host"
)

func TestDatabase_integration(t *testing.T) {
	var (
		l   = zaptest.NewLogger(t).Sugar()
		ctx = context.Background()
	)
	godotenv.Load("../.env")
	var installation = dev.GetTestInstallationID()
	if installation == "" {
		installation = "6969696969"
	}
	client, err := db.New(l.Named("db"), "integration_test", dev.DatabaseOptions)
	require.NoError(t, err)

	// make a repo
	var repos = client.Repos()
	err = repos.NewRepository(ctx, host.HostGitHub,
		installation, "bobheadxi", "calories")
	require.NoError(t, err)
	// get that repo id
	repo, err := repos.GetRepository(ctx, "bobheadxi", "calories")
	require.NoError(t, err)
	assert.NotZero(t, repo.ID)
	t.Log("bobheadxi/calories created as ID:", repo.ID)
	defer repos.DeleteRepository(ctx, repo.ID)

	// run tests
	t.Run("test host items", func(t *testing.T) {
		assert.NoError(t, repos.InsertHostItems(ctx, repo.ID, []*host.Item{
			{
				GitHubID: 1234,
				Number:   25,
				Type:     host.ItemTypeIssue,
				Details: map[string]interface{}{
					"some_detail": 23847125,
				},
			},
			// it is possible to have nil items pad the end
			nil,
			nil,
		}))
	})
	t.Run("test git burndown", func(t *testing.T) {
		assert.NoError(t, repos.InsertGitBurndownResult(ctx, repo.ID,
			testdata.Meta,
			testdata.Burndown))
		bd, err := repos.GetGlobalBurndown(ctx, repo.ID)
		require.NoError(t, err)
		assert.Equal(t, len(testdata.Burndown.Global), len(bd))
	})
}
