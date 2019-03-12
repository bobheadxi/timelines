package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const repoJobsName = "repojobs"

// RepoJob denotes a repository processing job
type RepoJob struct {
	ID uuid.UUID

	Owner          string
	Repo           string
	InstallationID string
}

// RepoJobState denotes the state of a job
type RepoJobState struct {
	Analysis   *StateMeta `json:"analysis"`
	GitHubSync *StateMeta `json:"github_sync"`
}

// RepoJobsClient exposes an API for interacting with repo job entries
type RepoJobsClient struct {
	c *Client
	l *zap.SugaredLogger
}

// Queue queues a repo for update
func (r *RepoJobsClient) Queue(job *RepoJob) error {
	var (
		l = r.l.With("job.id", job.ID)
	)

	// push job into queue
	b, err := json.Marshal(job)
	if err != nil {
		l.Errorw("failed to encode job",
			"error", err)
		return err
	}
	if err := r.c.redis.RPush(queueRepoJobs, string(b)).Err(); err != nil {
		l.Errorw("failed to queue job",
			"error", err)
		return err
	}

	// set job state tracker
	if err := r.SetState(job.ID, &RepoJobState{
		Analysis:   &StateMeta{State: StateNoProgress},
		GitHubSync: &StateMeta{State: StateNoProgress},
	}); err != nil {
		l.Errorw("failed to set default state for job",
			"error", err)
	}

	// execute pipe
	l.Infow("job queued")
	return nil
}

// Dequeue grabs the next repo job
func (r *RepoJobsClient) Dequeue(timeout time.Duration) (<-chan *RepoJob, <-chan error) {
	var (
		jobC = make(chan *RepoJob, 1)
		errC = make(chan error, 1)
	)

	go func() {
		defer close(errC)
		defer close(jobC)

		var pop = r.c.redis.BLPop(timeout, queueRepoJobs)
		data, err := pop.Result()
		if err == redis.Nil {
			jobC <- nil // indicate nothing was found
			return
		} else if err != nil {
			errC <- fmt.Errorf("error occured when popping from queue: %s", err.Error())
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
		l      = r.l.With("job.id", jobID)
		jobKey = stateKeyRepoJob(jobID)
		get    = r.c.redis.MGet(
			stateKeyRepoJobAnalysis(jobKey),
			stateKeyRepoJobGitHubSync(jobKey),
		)
	)

	data, err := get.Result()
	if err != nil {
		l.Warnw("could not find job state", "error", err)
		return nil, err
	}

	return &RepoJobState{
		Analysis:   ParseState(data[0]),
		GitHubSync: ParseState(data[1]),
	}, nil
}

// SetState updates the state of the given job ID with the given state
func (r *RepoJobsClient) SetState(jobID uuid.UUID, state *RepoJobState) error {
	var (
		l      = r.l.With("job.id", jobID)
		jobKey = stateKeyRepoJob(jobID)
		pipe   = r.c.redis.TxPipeline()
	)

	// set job state tracker - TODO: this is verbose
	if state.Analysis != nil {
		if state.Analysis.Meta == nil {
			state.Analysis.Meta = make(map[string]interface{})
		}
		state.Analysis.Meta["updated"] = time.Now()

		// add to pipeline
		data, _ := json.Marshal(state.Analysis)
		pipe.Set(
			stateKeyRepoJobAnalysis(jobKey),
			data,
			time.Hour)
	}
	if state.GitHubSync != nil {
		if state.GitHubSync.Meta == nil {
			state.GitHubSync.Meta = make(map[string]interface{})
		}
		state.GitHubSync.Meta["updated"] = time.Now()

		// add to pipeline
		data, _ := json.Marshal(state.GitHubSync)
		pipe.Set(
			stateKeyRepoJobGitHubSync(jobKey),
			data,
			time.Hour)
	}

	// execute pipe
	if _, err := pipe.Exec(); err != nil {
		l.Errorw("failed to update job",
			"error", err)
		return err
	}

	return nil
}
