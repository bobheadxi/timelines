package github

import (
	"context"

	"github.com/google/go-github/github"
)

type PullRequest struct {
	github.PullRequest
	Labels []github.Label
}

func NewPullRequest(ctx context.Context, cli *github.Client, i *github.Issue) (*github.PullRequest, error) {
	pr, _, err := cli.PullRequests.Get(ctx, "", "", *i.Number)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
