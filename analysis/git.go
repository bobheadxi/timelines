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
	pipe *hercules.Pipeline
	opts *GitRepoAnalyserOptions

	// repo metadata
	first time.Time
	size  int

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
	var (
		first time.Time
		size  int
	)
	history.ForEach(func(obj *object.Commit) error {
		// if no more parents, this is probably the first commit
		if obj.NumParents() == 0 {
			first = obj.Committer.When
		}
		size++
		return nil
	})

	return &GitRepoAnalyser{
		pipe: pipe,
		opts: &opts,

		first: first,
		size:  size,

		l: l,
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
	commits, err := g.pipe.Commits(true)
	if err != nil {
		return nil, err
	}
	g.pipe.Initialize(map[string]interface{}{
		leaves.ConfigBurndownTrackPeople: true,
	})
	return g.pipe.Run(commits)
}

func (g *GitRepoAnalyser) burndown() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.BurndownAnalysis{
		TrackFiles: true,

		// TODO: these is important to keep high (~30+) to account for largish (1k+)
		// to massive (100k+) repositories, but for smaller or shorter projects
		// (eg hackathons) this gives kind of useless data. will probably need some
		// way to scale this automatically, based on the repo size.
		Granularity: 30,
		Sampling:    30,

		PeopleNumber: 10,
	}).(hercules.LeafPipelineItem)
}

func (g *GitRepoAnalyser) coupling() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.CouplesAnalysis{
		PeopleNumber: 10,
	}).(hercules.LeafPipelineItem)
}
