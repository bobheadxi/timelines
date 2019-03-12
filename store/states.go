package store

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	statesRepoJobs = "states:" + repoJobsName + ":"
)

// StateMeta wraps State for additional metadata
type StateMeta struct {
	State   State
	Message string
	Meta    map[string]interface{}
}

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
func ParseState(v interface{}) *StateMeta {
	if v == nil {
		return nil
	}
	var (
		state = &StateMeta{}
		data  = v.(string)
	)
	json.Unmarshal([]byte(data), state)
	return state
}

func stateKeyRepoJob(jobID uuid.UUID) string         { return statesRepoJobs + jobID.String() }
func stateKeyRepoJobAnalysis(jobKey string) string   { return jobKey + ":analysis" }
func stateKeyRepoJobGitHubSync(jobKey string) string { return jobKey + ":github_sync" }
