package analysis

import "gopkg.in/src-d/hercules.v10/leaves"

// BurndownResult represents the burndown analysis result
type BurndownResult struct {
	Global          [][]int64            // [number of samples][number of bands]
	FileHistories   map[string][][]int64 // [file][number of samples][number of bands]
	PeopleHistories [][][]int64          // [people][number of samples][number of bands]

	FileOwnership map[string]map[int]int // [file][developer][lines]
	PeopleMatrix  [][]int64              // [number of people][number of people + 2] (map people -> lines)
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
