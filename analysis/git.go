package analysis

import (
	// TODO: explore gitbase
	// gitbase "gopkg.in/src-d/gitbase.v0"

	gogit "gopkg.in/src-d/go-git.v4"
	hercules "gopkg.in/src-d/hercules.v9"
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
