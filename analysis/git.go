package analysis

import (

	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	"time"

	"go.uber.org/zap"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	hercules "gopkg.in/src-d/hercules.v9"
	"gopkg.in/src-d/hercules.v9/leaves"
)

// GitRepoAnalyser executes pipelines on a repo
type GitRepoAnalyser struct {
	pipe  *hercules.Pipeline
	opts  *GitRepoAnalyserOptions
	first time.Time

	l *zap.SugaredLogger
}

// GitRepoAnalyserOptions denotes options for the analyzer
type GitRepoAnalyserOptions struct{}

// NewGitAnalyser sets up a new pipeline for repo analysis
func NewGitAnalyser(
	l *zap.SugaredLogger,
	repo *gogit.Repository,
	opts GitRepoAnalyserOptions,
) *GitRepoAnalyser {
	var pipe = hercules.NewPipeline(repo)
	pipe.PrintActions = false
	pipe.DumpPlan = false

	// get time of first commit, to calculate relative timeframes
	history, _ := repo.Log(&gogit.LogOptions{
		Order: gogit.LogOrderCommitterTime,
	})
	var first time.Time
	history.ForEach(func(obj *object.Commit) error {
		// if no more parents, this is probably the first commit
		if obj.NumParents() == 0 {
			first = obj.Committer.When
		}
		return nil
	})

	return &GitRepoAnalyser{
		pipe:  pipe,
		opts:  &opts,
		first: first,
	}
}

// GitRepoReport is a container around different analysis results
type GitRepoReport struct {
	First time.Time

	Burndown *BurndownResult
	Coupling *CouplingResult
}

// Analyze executes the pipeline
func (g *GitRepoAnalyser) Analyze() (*GitRepoReport, error) {

	// set up pipe
	var (
		burnItem     = g.burndown()
		couplingItem = g.coupling()
	)

	// execute analysis
	results, err := g.exec()
	if err != nil {
		return nil, err
	}

	// collect results
	var (
		burndown = newBurndownResult(results[burnItem].(leaves.BurndownResult))
		coupling = newCouplingResult(results[couplingItem].(leaves.CouplesResult))
	)

	return &GitRepoReport{
		First: g.first,

		Burndown: &burndown,
		Coupling: &coupling,
	}, nil
}

func (g *GitRepoAnalyser) exec() (map[hercules.LeafPipelineItem]interface{}, error) {
	commits, err := g.pipe.Commits(false)
	if err != nil {
		return nil, err
	}
	g.pipe.Initialize(nil)
	return g.pipe.Run(commits)
}

func (g *GitRepoAnalyser) burndown() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.BurndownAnalysis{
		TrackFiles:  true,
		Granularity: 30,
		Sampling:    30,

		PeopleNumber: 10, // TODO: this should scale with actual contributors
	}).(hercules.LeafPipelineItem)
}

func (g *GitRepoAnalyser) coupling() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.CouplesAnalysis{
		PeopleNumber: 10,
	}).(hercules.LeafPipelineItem)
}
