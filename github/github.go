package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Client handles GitHub sync and API requests
type Client struct {
	gh *github.Client

	l *zap.SugaredLogger
}

// NewClient instantiates a new github client
func NewClient(ctx context.Context, l *zap.SugaredLogger, token string) (*Client, error) {
	var tc *http.Client
	if token != "" {
		l.Infow("loading oauth client")
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	} else {
		l.Infow("using default client")
		tc = http.DefaultClient
	}

	var gh = github.NewClient(tc)
	lim, _, err := gh.RateLimits(ctx)
	if err != nil {
		return nil, err
	}
	l.Infow("authenticated against github",
		"rate_limits", lim)

	return &Client{gh: gh}, nil
}

// IssueState denotes whether to get issues that are closed, open, or all
type IssueState string

const (
	// IssueStateAll = all issues
	IssueStateAll IssueState = "all"
	// IssueStateClosed = closed issues
	IssueStateClosed IssueState = "closed"
	// IssueStateOpen = open issues
	IssueStateOpen IssueState = "open"
)

// IssuesFilter denotes options to filter issues by
type IssuesFilter struct {
	MinIssue int
	State    IssueState
	Interval time.Duration
}

// GetIssues retrieves all issues for a project
func (c *Client) GetIssues(
	ctx context.Context,
	user, repo string,
	filter IssuesFilter,
	issuesC chan<- *github.Issue,
	prsC chan<- *github.Issue,
) error {
	var (
		l         = c.l.With("user", user, "repo", repo)
		itemCount = 0
	)
	for page := filter.MinIssue/100 + 1; page != 0; {
		items, resp, err := c.gh.Issues.ListByRepo(ctx, user, repo, &github.IssueListByRepoOptions{
			Direction: "asc",
			Sort:      "created",
			State:     string(filter.State),
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			l.Errorw("failed to fetch issues", "error", err, "items", items, "page", page)
			return fmt.Errorf("failed to fetch issues for '%s/%s' on page '%d': %s",
				user, repo, page, err)
		}

		itemCount += len(items)
		l.Infow("items retrieved", "items", itemCount, "page", page)
		for _, i := range items {
			if !i.IsPullRequest() {
				issuesC <- i
			} else {
				prsC <- i
			}
		}

		page = resp.NextPage
		if filter.Interval > 0 {
			time.Sleep(filter.Interval)
		}
	}

	return nil
}

// GetPullRequest extracts a Pull Request from the given issue
func (c *Client) GetPullRequest(ctx context.Context, i *github.Issue) (*github.PullRequest, error) {
	if !i.IsPullRequest() {
		return nil, fmt.Errorf("issue '%d' is not a pull request", i.GetNumber())
	}

	var repo = i.GetRepository()
	pr, _, err := c.gh.PullRequests.Get(ctx,
		repo.GetOwner().GetName(),
		repo.GetName(),
		i.GetNumber())
	if err != nil {
		c.l.Errorw("failed to get pull request",
			"issue", i.GetNumber(),
			"repo", repo)
		return nil, fmt.Errorf("failed to get pull request for issue '%d'", i.GetNumber())
	}

	return pr, nil
}
