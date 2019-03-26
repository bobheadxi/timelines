package analysis

import (

	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	"time"

	"go.uber.org/zap"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	hercules "gopkg.in/src-d/hercules.v10"
	"gopkg.in/src-d/hercules.v10/leaves"
)

// GitRepoAnalyser executes pipelines on a repo
type GitRepoAnalyser struct {
	pipe *hercules.Pipeline
	opts *GitRepoAnalyserOptions

	// repo metadata
	first time.Time
	last  time.Time
	size  int

	l *zap.SugaredLogger

	// analysis metadata
	facts map[string]interface{}
}

// GitRepoAnalyserOptions denotes options for the analyzer
type GitRepoAnalyserOptions struct {
	TickCount int
}

// NewGitAnalyser sets up a new pipeline for repo analysis
func NewGitAnalyser(
	l *zap.SugaredLogger,
	repo *gogit.Repository,
	opts GitRepoAnalyserOptions,
) (*GitRepoAnalyser, error) {
	var pipe = hercules.NewPipeline(repo)
	pipe.PrintActions = false
	pipe.DumpPlan = false
	if opts.TickCount == 0 {
		opts.TickCount = 500
	}

	// get time of first commit, to calculate relative timeframes
	history, err := repo.Log(&gogit.LogOptions{
		Order: gogit.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}
	var (
		first, last time.Time
		size        int
	)
	obj, err := history.Next()
	if err != nil {
		return nil, err
	}
	last = obj.Committer.When
	history.ForEach(func(obj *object.Commit) error {
		// if no more parents, this is probably the first commit
		if obj.NumParents() == 0 {
			first = obj.Committer.When
		}
		size++
		return nil
	})
	l.Infow("repo prepped",
		"size", size,
		"first_commit", first,
		"last_commit", last)

	return &GitRepoAnalyser{
		pipe: pipe,
		opts: &opts,

		first: first,
		last:  last,
		size:  size,

		l: l,
	}, nil
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
		people   = g.facts[hercules.FactIdentityDetectorReversedPeopleDict].([]string)
		burndown = newBurndownResult(results[burnItem].(leaves.BurndownResult), people)
		coupling = newCouplingResult(results[couplingItem].(leaves.CouplesResult))
	)

	return &GitRepoReport{
		First: g.first,

		Burndown: &burndown,
		Coupling: &coupling,
	}, nil
}

func (g *GitRepoAnalyser) exec() (map[hercules.LeafPipelineItem]interface{}, error) {
	// calculate number of ticks to use
	hours := g.last.Sub(g.first).Hours()
	tickSize := int(hours / float64(g.opts.TickCount))
	g.l.Infow("preparing pipeline",
		"total_hours", hours,
		"tick_size", tickSize)

	// grab commits and initialize pipeline
	commits, err := g.pipe.Commits(true)
	if err != nil {
		return nil, err
	}
	g.facts = map[string]interface{}{
		"TicksSinceStart.TickSize": tickSize,

		// Tree config
		"TreeDiff.Languages":           []string{"all"},
		"TreeDiff.EnableBlacklist":     true,
		"TreeDiff.BlacklistedPrefixes": []string{"vendor/", "vendors/", "node_modules/"},

		// Burndown config
		leaves.ConfigBurndownTrackPeople: true,
		leaves.ConfigBurndownTrackFiles:  true,
	}
	g.pipe.Initialize(g.facts)

	// execute pipeline
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
		Sampling:    30, // sampling != granularity seems broken - see https://github.com/src-d/hercules/issues/260

		PeopleNumber: 10,
	}).(hercules.LeafPipelineItem)
}

func (g *GitRepoAnalyser) coupling() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.CouplesAnalysis{
		PeopleNumber: 10,
	}).(hercules.LeafPipelineItem)
}
