package store

import (
	"os"
	"testing"
	"time"

	"github.com/bobheadxi/timelines/dev"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestStore(t *testing.T) {
	godotenv.Load("../.env")

	l := zaptest.NewLogger(t).Sugar()
	c, err := NewClient(l, dev.StoreOptions)
	assert.NoError(t, err)
	defer c.Reset()

	// add job
	id, _ := uuid.NewUUID()
	err = c.RepoJobs().Queue(&RepoJob{
		ID:             id,
		Owner:          "bobheadxi",
		Repo:           "calories",
		InstallationID: os.Getenv("GITHUB_TEST_INSTALLTION"),
	})
	assert.NoError(t, err)

	// get job
	jobC, errC := c.RepoJobs().Dequeue(5 * time.Second)
	select {
	case job := <-jobC:
		assert.Equal(t, id, job.ID)
		assert.Equal(t, "bobheadxi", job.Owner)
	case err := <-errC:
		assert.NoError(t, err)
		t.Fail()
	}

	// get the job state
	state, err := c.RepoJobs().GetState(id)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	t.Log(state)
	assert.Equal(t, StateNoProgress, state.Analysis.State)

	// update job state
	err = c.RepoJobs().SetState(id, &RepoJobState{
		Analysis: &StateMeta{
			State:   StateDone,
			Message: "hello world",
		},
	})
	assert.NoError(t, err)

	// get updated job state
	state, err = c.RepoJobs().GetState(id)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	t.Log(state)
	assert.Equal(t, StateDone, state.Analysis.State)
	assert.Equal(t, "hello world", state.Analysis.Message)
	assert.Equal(t, StateNoProgress, state.GitHubSync.State)
}
