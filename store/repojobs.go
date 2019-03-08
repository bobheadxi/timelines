package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const repoJobsName = "repojobs"

// RepoJob denotes a repository processing job
type RepoJob struct {
	ID uuid.UUID

	Owner string
	Repo  string
}

// RepoJobState denotes the state of a job
type RepoJobState struct {
	Analysis   State
	GitHubSync State
}

// RepoJobsClient exposes an API for interacting with repo job entries
type RepoJobsClient struct {
	c *Client
	l *zap.SugaredLogger
}

// Queue queues a repo for update
func (r *RepoJobsClient) Queue(job *RepoJob) error {
	var (
		l    = r.l.With("job.id", job.ID)
		b    bytes.Buffer
		enc  = json.NewEncoder(&b)
		pipe = r.c.redis.TxPipeline()
	)

	// push job into queue
	if err := enc.Encode(job); err != nil {
		l.Errorw("failed to encode job",
			"error", err)
		return err
	}
	pipe.RPush(queueRepoJobs, b.String())

	// set job state tracker
	b.Reset()
	if err := enc.Encode(&RepoJobState{}); err != nil {
		l.Errorw("failed to encode job",
			"error", err)
		return err
	}
	pipe.Set(stateKeyRepoJob(job.ID), b.String(), time.Hour)

	// execute pipe
	if _, err := pipe.Exec(); err != nil {
		l.Errorw("failed to push job",
			"error", err)
		return err
	}
	l.Infow("job queued")

	return nil
}

// Dequeue grabs the next repo job
func (r *RepoJobsClient) Dequeue() (<-chan *RepoJob, <-chan error) {
	var (
		jobC = make(chan *RepoJob, 1)
		errC = make(chan error, 1)
	)

	go func() {
		defer close(errC)
		defer close(jobC)

		var pop = r.c.redis.BLPop(time.Second, queueRepoJobs)
		data, err := pop.Result()
		if err != nil {
			errC <- err
			return
		} else if len(data) < 2 {
			errC <- errors.New("nothing was popped from queue")
			return
		}

		var job = &RepoJob{}
		if err := json.Unmarshal([]byte(data[1]), job); err != nil {
			r.l.Errorw("failed to read data", "error", err)
			errC <- err
			return
		}

		r.l.Infow("job dequeued", "job.id", job.ID)
		jobC <- job
	}()

	return jobC, errC
}

// GetState retrieves the state of the given job ID
func (r *RepoJobsClient) GetState(jobID uuid.UUID) (*RepoJobState, error) {
	var (
		l   = r.l.With("job.id", jobID)
		get = r.c.redis.Get(stateKeyRepoJob(jobID))
	)

	data, err := get.Bytes()
	if err != nil {
		return nil, err
	}

	var state = &RepoJobState{}
	if err := json.Unmarshal(data, state); err != nil {
		l.Errorw("failed to read data", "error", err)
		return nil, err
	}

	return state, nil
}

// SetState updates the state of the given job ID with the given state
func (r *RepoJobsClient) SetState(jobID uuid.UUID, state *RepoJobState) error {
	var (
		l   = r.l.With("job.id", jobID)
		b   bytes.Buffer
		enc = json.NewEncoder(&b)
	)

	// encode data
	if err := enc.Encode(state); err != nil {
		l.Errorw("failed to encode job",
			"error", err)
		return err
	}

	// set job state tracker
	var set = r.c.redis.Set(stateKeyRepoJob(jobID), b.String(), time.Hour)
	if err := set.Err(); err != nil {
		l.Errorw("failed to update job",
			"error", err)
		return err
	}

	return nil
}
