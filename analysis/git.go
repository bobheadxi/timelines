package analysis

import (
	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	// TODO: revert back to src-d/hercules.v9
	// https://github.com/src-d/hercules/pull/230
	"github.com/bobheadxi/hercules"

	gogit "gopkg.in/src-d/go-git.v4"
)

type GitRepoAnalyser struct {
	pipe *hercules.Pipeline
}

func NewGit(repo *gogit.Repository) *GitRepoAnalyser {
	var pipe = hercules.NewPipeline(repo)
	pipe.PrintActions = false

	return &GitRepoAnalyser{
		pipe: pipe,
	}
}
