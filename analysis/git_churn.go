package analysis

import (
	"io"
	"sort"
	"unicode/utf8"

	"github.com/sergi/go-diff/diffmatchpatch"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	hercules "gopkg.in/src-d/hercules.v9"
)

// ChurnAnalysisResult represents the churn analysis result
type ChurnAnalysisResult struct {
	Global Edits
	People map[string]Edits
}

// Modified from https://github.com/src-d/hercules/blob/master/contrib/_plugin_example/churn_analysis.go
// originally authored by @vmarkovtsev
type churnAnalysis struct {
	// No special merge logic is required
	hercules.NoopMerger
	// Process each merge only once
	hercules.OneShotMergeProcessor

	global             []editInfo
	people             map[int][]editInfo
	reversedPeopleDict []string
}

// Edits denotes edit data
type Edits struct {
	Days      []int
	Additions []int
	Removals  []int
}

// Requires returns the list of dependencies which must be supplied in Consume().
// file_diff - line diff for each commit change
// changes - list of changed files for each commit
// blob_cache - set of blobs affected by each commit
// day - number of days since start for each commit
// author - author of the commit
func (c *churnAnalysis) Requires() []string {
	arr := [...]string{
		hercules.DependencyFileDiff,
		hercules.DependencyTreeChanges,
		hercules.DependencyBlobCache,
		hercules.DependencyDay,
		hercules.DependencyAuthor}
	return arr[:]
}

// Configure applies the parameters specified in the command line. Map keys correspond to "Name".
func (c *churnAnalysis) Configure(facts map[string]interface{}) error {
	c.reversedPeopleDict = facts[hercules.FactIdentityDetectorReversedPeopleDict].([]string)
	return nil
}

// Initialize resets the internal temporary data structures and prepares the object for Consume().
func (c *churnAnalysis) Initialize(repository *gogit.Repository) error {
	c.global = []editInfo{}
	c.people = map[int][]editInfo{}
	c.OneShotMergeProcessor.Initialize()
	return nil
}

func (c *churnAnalysis) Consume(deps map[string]interface{}) (map[string]interface{}, error) {
	if !c.ShouldConsumeCommit(deps) {
		return nil, nil
	}
	fileDiffs := deps[hercules.DependencyFileDiff].(map[string]hercules.FileDiffData)
	treeDiffs := deps[hercules.DependencyTreeChanges].(object.Changes)
	cache := deps[hercules.DependencyBlobCache].(map[plumbing.Hash]*hercules.CachedBlob)
	day := deps[hercules.DependencyDay].(int)
	author := deps[hercules.DependencyAuthor].(int)
	for _, change := range treeDiffs {
		action, err := change.Action()
		if err != nil {
			return nil, err
		}
		added := 0
		removed := 0
		switch action {
		case merkletrie.Insert:
			added, _ = cache[change.To.TreeEntry.Hash].CountLines()
		case merkletrie.Delete:
			removed, _ = cache[change.From.TreeEntry.Hash].CountLines()
		case merkletrie.Modify:
			diffs := fileDiffs[change.To.Name]
			for _, edit := range diffs.Diffs {
				length := utf8.RuneCountInString(edit.Text)
				switch edit.Type {
				case diffmatchpatch.DiffEqual:
					continue
				case diffmatchpatch.DiffInsert:
					added += length
				case diffmatchpatch.DiffDelete:
					removed += length
				}
			}

		}
		if err != nil {
			return nil, err
		}
		ei := editInfo{Day: day, Added: added, Removed: removed}
		c.global = append(c.global, ei)
		seq, exists := c.people[author]
		if !exists {
			seq = []editInfo{}
		}
		seq = append(seq, ei)
		c.people[author] = seq
	}
	return nil, nil
}

func (c *churnAnalysis) Finalize() interface{} {
	result := ChurnAnalysisResult{
		Global: editInfosToEdits(c.global),
		People: map[string]Edits{},
	}
	for key, val := range c.people {
		result.People[c.reversedPeopleDict[key]] = editInfosToEdits(val)
	}
	return result
}

// Functions to fulfill hercules.LeafPipeline

func (c *churnAnalysis) Name() string                                 { return "churnAnalysis" }
func (c *churnAnalysis) Provides() []string                           { return []string{} }
func (c *churnAnalysis) Flag() string                                 { return "" }
func (c *churnAnalysis) Description() string                          { return "" }
func (c *churnAnalysis) Serialize(interface{}, bool, io.Writer) error { return nil }
func (c *churnAnalysis) ListConfigurationOptions() []hercules.ConfigurationOption {
	return []hercules.ConfigurationOption{}
}
func (c *churnAnalysis) Fork(n int) []hercules.PipelineItem {
	return hercules.ForkSamePipelineItem(c, n)
}

// Utilities

type editInfo struct {
	Day     int
	Added   int
	Removed int
}

func editInfosToEdits(eis []editInfo) Edits {
	aux := map[int]*editInfo{}
	for _, ei := range eis {
		ptr := aux[ei.Day]
		if ptr == nil {
			ptr = &editInfo{Day: ei.Day}
		}
		ptr.Added += ei.Added
		ptr.Removed += ei.Removed
		aux[ei.Day] = ptr
	}
	seq := []int{}
	for key := range aux {
		seq = append(seq, key)
	}
	sort.Ints(seq)
	edits := Edits{
		Days:      make([]int, len(seq)),
		Additions: make([]int, len(seq)),
		Removals:  make([]int, len(seq)),
	}
	for i, day := range seq {
		edits.Days[i] = day
		edits.Additions[i] = aux[day].Added
		edits.Removals[i] = aux[day].Removed
	}
	return edits
}
