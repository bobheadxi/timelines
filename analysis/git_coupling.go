package analysis

import "gopkg.in/src-d/hercules.v9/leaves"

// CouplingResult represents the coupling analysis result
type CouplingResult struct {
	FilesMatrix []map[int]int64
	FilesLines  []int
	Files       []string
}

func newCouplingResult(r leaves.CouplesResult) CouplingResult {
	return CouplingResult{
		FilesMatrix: r.FilesMatrix,
		FilesLines:  r.FilesLines,
		Files:       r.Files,
	}
}
