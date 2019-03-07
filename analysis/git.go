package analysis

import (

	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	gogit "gopkg.in/src-d/go-git.v4"
	hercules "gopkg.in/src-d/hercules.v9"
	"gopkg.in/src-d/hercules.v9/leaves"
)

// GitRepoAnalyser executes pipelines on a repo
type GitRepoAnalyser struct {
	pipe *hercules.Pipeline
}

// NewGitAnalyser sets up a new pipeline for repo analysis
func NewGitAnalyser(repo *gogit.Repository) *GitRepoAnalyser {
	var pipe = hercules.NewPipeline(repo)
	pipe.PrintActions = false
	pipe.DumpPlan = false

	return &GitRepoAnalyser{
		pipe: pipe,
	}
}

// GitRepoReport is a container around different analysis results
type GitRepoReport struct {
	Burndown *BurndownResult
	Churn    *ChurnAnalysisResult
}

// Analyze executes the pipeline
func (g *GitRepoAnalyser) Analyze() (*GitRepoReport, error) {

	// set up pipe
	var (
		burnItem  = g.burndown()
		churnItem = g.churn()
	)

	// execute analysis
	results, err := g.exec()
	if err != nil {
		return nil, err
	}

	// collect results
	var (
		burndown = newBurndownResult(results[burnItem].(leaves.BurndownResult))
		churn    = results[churnItem].(ChurnAnalysisResult)
	)

	return &GitRepoReport{
		Burndown: &burndown,
		Churn:    &churn,
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

func (g *GitRepoAnalyser) churn() hercules.LeafPipelineItem {
	return g.pipe.DeployItem(&churnAnalysis{}).(hercules.LeafPipelineItem)
}
