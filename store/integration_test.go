package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestStore(t *testing.T) {
	l := zaptest.NewLogger(t).Sugar()
	c, err := NewClient(l, devEnvOptions)
	assert.NoError(t, err)
	defer c.redis.Del(queueRepoJobs, statesRepoJobs)

	// add job
	id, _ := uuid.NewUUID()
	err = c.RepoJobs().Queue(&RepoJob{
		ID:    id,
		Owner: "bobheadxi",
		Repo:  "calories",
	})
	assert.NoError(t, err)

	// get job
	jobC, errC := c.RepoJobs().Dequeue()
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
	assert.Equal(t, StateNoProgress, state.Analysis)

	// update job state
	err = c.RepoJobs().SetState(id, &RepoJobState{
		Analysis: StateDone,
	})
	assert.NoError(t, err)

	// get updated job state
	state, err = c.RepoJobs().GetState(id)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	t.Log(state)
	assert.Equal(t, StateDone, state.Analysis)
}
