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
	gh    *github.Client
	token string

	l *zap.SugaredLogger
}

func newClient(l *zap.SugaredLogger, auth oauth2.TokenSource) (*Client, error) {
	var (
		tc    *http.Client
		token string
	)
	if auth != nil {
		l.Info("loading credentials")
		t, err := auth.Token()
		if err != nil {
			l.Errorw("failed to load credentials", "error", err)
			return nil, err
		}
		token = t.AccessToken
		tc = oauth2.NewClient(oauth2.NoContext, oauth2.ReuseTokenSource(t, auth))
	} else {
		l.Infow("using default client")
		tc = http.DefaultClient
	}

	return &Client{
		gh:    github.NewClient(tc),
		token: token,
		l:     l,
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

// InstallationToken returns the token in use by the client, if there is one
func (c *Client) InstallationToken() string { return c.token }

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
	wait *sync.WaitGroup,
) error {
	var (
		l         = c.l.With("user", user, "repo", repo, "sync", "issues")
		itemCount = 0
	)

	defer func() {
		close(issuesC)
		l.Infof("collected %d items", itemCount)
		wait.Done()
	}()

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
		l.Infow("items retrieved", "total_items", itemCount, "page", page)

		// queue all items into output
		for _, i := range items {
			if !i.IsPullRequest() {
				issuesC <- i
			}
			// TODO: some metadata is only available on this API, and not the PR one,
			// most notably reactions. somehow merge results with the PR fetching?
		}

		// wait if required before making next request
		page = resp.NextPage
		if filter.Interval > 0 {
			time.Sleep(filter.Interval)
		}
	}

	return nil
}

// GetPullRequests retrieves all pull requests for a project
func (c *Client) GetPullRequests(
	ctx context.Context,
	user, repo string,
	filter ItemFilter,
	prC chan<- *github.PullRequest,
	wait *sync.WaitGroup,
) error {
	var (
		l         = c.l.With("user", user, "repo", repo, "sync", "pull_requests")
		itemCount = 0
	)

	defer func() {
		close(prC)
		l.Infof("collected %d items", itemCount)
		wait.Done()
	}()

	for page := filter.MinNumber/100 + 1; page != 0; {

		// fetch items
		items, resp, err := c.gh.PullRequests.List(ctx, user, repo, &github.PullRequestListOptions{
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
		l.Infow("items retrieved", "total_items", itemCount, "page", page)

		// queue all items into output
		for _, pr := range items {
			prC <- pr
		}

		// wait if required before making next request
		page = resp.NextPage
		if filter.Interval > 0 {
			time.Sleep(filter.Interval)
		}
	}

	return nil
}
