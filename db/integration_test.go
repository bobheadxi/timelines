package db_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/dev"
	"github.com/bobheadxi/timelines/github"
)

func TestDatabase(t *testing.T) {
	var (
		l   = zaptest.NewLogger(t).Sugar()
		ctx = context.Background()
	)
	godotenv.Load("../.env")
	var installation = dev.GetTestInstallationID()
	if installation == "" {
		installation = "6969696969"
	}
	client, err := db.New(l, "integration_test", dev.DatabaseOptions)
	assert.NoError(t, err)

	var repos = client.Repos()

	// make a repo
	err = repos.NewRepository(ctx, installation, "bobheadxi", "calories")
	assert.NoError(t, err)

	// get that repo id
	id, err := repos.GetRepositoryID(ctx, "bobheadxi", "calories")
	assert.NoError(t, err)
	assert.NotZero(t, id)
	t.Log("bobheadxi/calories created as ID:", id)

	// add an item
	err = repos.InsertGitHubItems(ctx, id, []*github.Item{
		&github.Item{
			GitHubID: 1234,
			Number:   25,
			Type:     github.ItemTypeIssue,
			Details: map[string]interface{}{
				"some_detail": 23847125,
			},
		},
	})
	assert.NoError(t, err)

	// delete repo
	err = repos.DeleteRepository(ctx, id)
	assert.NoError(t, err)
}
