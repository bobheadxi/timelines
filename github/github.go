package github

import (
	"context"
	"fmt"
	"net/http"
	"sync"
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

func newClient(l *zap.SugaredLogger, auth oauth2.TokenSource) (*Client, error) {
	var tc *http.Client
	if auth != nil {
		l.Info("loading credentials")
		token, err := auth.Token()
		if err != nil {
			l.Errorw("failed to load credentials", "error", err)
			return nil, err
		}
		l.Infow("token generated", "token", token)
		tc = oauth2.NewClient(oauth2.NoContext, oauth2.ReuseTokenSource(token, auth))
	} else {
		l.Infow("using default client")
		tc = http.DefaultClient
	}

	return &Client{
		gh: github.NewClient(tc),
		l:  l,
	}, nil
}

// NewClient instantiates a new github client. This package offers several
// implementations of auth, and different clients should be instantiated for
// different purposes.
//
// * a "signing client" should use AppAuth and NewSigningClient, and only be
//   used to create installation clients.
// * a "installation client" should use InstallationAuth, and only be created
//   from a SigningClient
// * a "default client" should use nil auth, and be used for unauthenticated
//   requests
//
func NewClient(ctx context.Context, l *zap.SugaredLogger, auth oauth2.TokenSource) (*Client, error) {
	client, err := newClient(l, auth)
	if err != nil {
		return nil, err
	}

	lim, _, err := client.gh.RateLimits(ctx)
	if err != nil {
		l.Infow("could not authenticate against github", "error", err)
		return nil, err
	}
	l.Infow("authenticated against github",
		"rate_limits", lim)

	return client, nil
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

// ItemFilter denotes options to filter issues and PRs by
type ItemFilter struct {
	MinNumber int
	State     IssueState
	Interval  time.Duration
}

// GetIssues retrieves all issues for a project
func (c *Client) GetIssues(
	ctx context.Context,
	user, repo string,
	filter ItemFilter,
	issuesC chan<- *github.Issue,
	fetchDetailsC chan<- *github.Issue,
	wait *sync.WaitGroup,
) error {
	wait.Add(1)
	defer wait.Done()

	var (
		l         = c.l.With("user", user, "repo", repo)
		itemCount = 0
	)

	for page := filter.MinNumber/100 + 1; page != 0; {

		// fetch items
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

		// queue all items into output
		for _, i := range items {
			if !i.IsPullRequest() {
				issuesC <- i
			} else {
				// set repository data to help fetch pull request details
				if i.GetRepository() == nil || i.GetRepository().GetOwner() == nil {
					i.Repository = &github.Repository{
						Name: github.String(repo),
						Owner: &github.User{
							Login: github.String(user),
						},
					}
				}
				fetchDetailsC <- i
			}
		}

		// wait if required before making next request
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
	if repo == nil || repo.GetOwner() == nil {
		return nil, fmt.Errorf("repo or owner is not set in issue '%d'", i.GetNumber())
	}
	pr, _, err := c.gh.PullRequests.Get(ctx,
		repo.GetOwner().GetLogin(),
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
