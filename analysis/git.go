package analysis

import (

	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	"errors"
	"time"

	"github.com/bobheadxi/timelines/log"
	"go.uber.org/zap"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	hercules "gopkg.in/src-d/hercules.v10"
	"gopkg.in/src-d/hercules.v10/leaves"
)

// GitRepoReport is a container around different analysis results
type GitRepoReport struct {
	Meta     *GitRepoMeta
	Burndown *BurndownResult
	Coupling *CouplingResult
}

// GitRepoMeta denotes the configuration used to execute the analysis
type GitRepoMeta struct {
	Commits  int
	First    time.Time
	Last     time.Time
	TickSize int
}

// GetTickRange returns the beginning and end times of the given tick
func (m *GitRepoMeta) GetTickRange(t int) (time.Time, time.Time) {
	real := time.Duration(m.TickSize) * time.Hour
	return m.First.Add(real * time.Duration(t)),
		m.First.Add(real * (time.Duration(t) + 1))
}

// GitRepoAnalyser executes pipelines on a repo
type GitRepoAnalyser struct {
	pipe *hercules.Pipeline
	l    *zap.SugaredLogger
	opts *GitRepoAnalyserOptions

	// repo metadata
	first    time.Time
	last     time.Time
	commits  int
	tickSize int

	// analysis metadata - this is populated by hercules
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
		opts.TickCount = 600
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
		if obj.NumParents() == 0 {
			// if no more parents, this is probably the first commit
			first = obj.Committer.When
		}
		size++
		return nil
	})
	if first.Unix() == 0 {
		return nil, errors.New("could not find a 'first' commit with no parents")
	}

	// calculate number of ticks to use
	hours := last.Sub(first).Hours()
	tickSize := int(hours / float64(opts.TickCount))
	if tickSize == 0 {
		tickSize = 1
	}
	l.Infow("repo prepped",
		"commits", size,
		"first_commit", first,
		"last_commit", last,
		"total_hours", hours,
		"tick_size", tickSize)

	return &GitRepoAnalyser{
		pipe: pipe,
		opts: &opts,
		l:    l,

		first:    first,
		last:     last,
		commits:  size,
		tickSize: tickSize,
	}, nil
}

// Analyze executes the pipeline
func (g *GitRepoAnalyser) Analyze() (*GitRepoReport, error) {

	// set up pipe
	var (
		start        = time.Now()
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

	// log some stuff about the results
	g.l.Infow("analysis complete",
		"duration", time.Since(start),
		"burndown.intervals", len(burndown.Global),
		"burndown.bands", len(burndown.Global[0]))

	return &GitRepoReport{
		Meta: &GitRepoMeta{
			Commits:  g.commits,
			First:    g.first,
			Last:     g.last,
			TickSize: g.tickSize,
		},
		Burndown: &burndown,
		Coupling: &coupling,
	}, nil
}

func (g *GitRepoAnalyser) exec() (map[hercules.LeafPipelineItem]interface{}, error) {
	// grab commits and initialize pipeline - only grab parents if there are over
	// 5000 commits for performance. TODO: assess threshold
	commits, err := g.pipe.Commits(g.commits > 5000)
	if err != nil {
		return nil, err
	}
	g.facts = map[string]interface{}{
		hercules.ConfigLogger:      log.NewHerculesLogger(g.l.Named("hercules")),
		"TicksSinceStart.TickSize": g.tickSize,

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
		TrackFiles:  true,
		Granularity: 30,
		Sampling:    30, // sampling != granularity seems broken - see https://github.com/src-d/hercules/issues/260
	}).(hercules.LeafPipelineItem)
}

func (g *GitRepoAnalyser) coupling() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&leaves.CouplesAnalysis{}).(hercules.LeafPipelineItem)
}
