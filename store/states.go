package store

import "github.com/google/uuid"

const (
	statesRepoJobs = "states:" + repoJobsName + ":"
)

// State denotes the state of a job
type State int

const (
	// StateNoProgress indicates nothing has happened yet
	StateNoProgress State = iota

	// StateInProgress indicates that the job is in progress
	StateInProgress

	// StateDone indicates the job was completed
	StateDone

	// StateError indicates something went wrong
	StateError
)

func stateKeyRepoJob(jobID uuid.UUID) string { return statesRepoJobs + jobID.String() }
