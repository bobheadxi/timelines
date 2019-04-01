package analysis

import "gopkg.in/src-d/hercules.v10/leaves"

// BurndownResult represents the burndown analysis result
type BurndownResult struct {
	// Burndowns are matrices that represent each band's values at a particular
	// sampling size, which will be 30 * ticks (at least for now - see the
	// analysis class for more details). From what I understand, it's sort of like
	// this:
	//     time intervals (aka ticks) * bands (representing if code is still present)
	// Might be able to drop Global - TODO: assess if sum of files per band and tick
	// is equal to global.
	Global [][]int64
	People map[string][][]int64
	Files  map[string][][]int64

	// Misc analysis
	// FileOwnership map[string]map[string]int
}

func newBurndownResult(r leaves.BurndownResult, people []string) BurndownResult {
	return BurndownResult{
		Global: r.GlobalHistory,
		People: peopleBurndowns(r.PeopleHistories, people),
		Files:  r.FileHistories,

		// FileOwnership: TODO
	}
}

func peopleBurndowns(data [][][]int64, people []string) map[string][][]int64 {
	res := make(map[string][][]int64)
	for i, bd := range data {
		if i > len(people) {
			continue // shouldn't happen, but just in case
		}
		res[people[i]] = bd
	}
	return res
}
