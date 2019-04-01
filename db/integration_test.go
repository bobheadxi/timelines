package db_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	// make a repo
	var repos = client.Repos()
	err = repos.NewRepository(ctx, host.HostGitHub,
		installation, "bobheadxi", "calories")
	assert.NoError(t, err)
	// get that repo id
	id, err := repos.GetRepositoryID(ctx, "bobheadxi", "calories")
	assert.NoError(t, err)
	assert.NotZero(t, id)
	t.Log("bobheadxi/calories created as ID:", id)
	defer repos.DeleteRepository(ctx, id)

	// run tests
	t.Run("test host items", func(t *testing.T) {
		assert.NoError(t, repos.InsertHostItems(ctx, id, []*host.Item{
			&host.Item{
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
		assert.NoError(t, repos.InsertGitBurndownResult(ctx, id,
			testdata.Meta,
			testdata.Burndown))
	})
}
