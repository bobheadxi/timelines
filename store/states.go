package store

import (
	"strconv"

	"github.com/google/uuid"
)

const (
	statesRepoJobs = "states:" + repoJobsName + ":"
)

// State denotes the state of a job
type State int

const (
	// StateNoProgress indicates nothing has happened yet
	StateNoProgress State = iota + 1

	// StateInProgress indicates that the job is in progress
	StateInProgress

	// StateDone indicates the job was completed
	StateDone

	// StateError indicates something went wrong
	StateError
)

// ParseState casts given interface to a State
func ParseState(v interface{}) State {
	if v == nil {
		return State(0)
	}
	s, _ := v.(string)
	i, _ := strconv.Atoi(s)
	return State(i)
}

func stateKeyRepoJob(jobID uuid.UUID) string         { return statesRepoJobs + jobID.String() }
func stateKeyRepoJobAnalysis(jobKey string) string   { return jobKey + ":analysis" }
func stateKeyRepoJobGitHubSync(jobKey string) string { return jobKey + ":github_sync" }
