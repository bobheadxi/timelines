package analysis

import "gopkg.in/src-d/hercules.v9/leaves"

// BurndownResult represents the burndown analysis result
type BurndownResult struct {
	Global          [][]int64
	FileHistories   map[string][][]int64
	FileOwnership   map[string]map[int]int
	PeopleHistories [][][]int64
	PeopleMatrix    [][]int64
}

func newBurndownResult(r leaves.BurndownResult) BurndownResult {
	return BurndownResult{
		Global:          r.GlobalHistory,
		FileHistories:   r.FileHistories,
		FileOwnership:   r.FileOwnership,
		PeopleHistories: r.PeopleHistories,
		PeopleMatrix:    r.PeopleMatrix,
	}
}
