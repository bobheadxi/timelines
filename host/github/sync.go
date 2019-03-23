package github

import (
	"context"
	"errors"
	"sync"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/host"
)

// SyncOptions denotes options for a syncer
type SyncOptions struct {
	Repo                Repo
	Filter              ItemFilter
	DetailsFetchWorkers int
	OutputBufferSize    int
}

// Repo denotes a repository to sync
type Repo struct {
	Owner string
	Name  string
}

// Syncer manages all GitHub synchronization tasks
type Syncer struct {
	c *Client
	l *zap.SugaredLogger

	opts SyncOptions

	issuesC chan *github.Issue
	prC     chan *github.PullRequest

	outC chan *host.Item

	used bool
}

// NewSyncer instantiates a new GitHub Syncer. It Syncer::Sync() can only be
// used once.
// TODO: should it be reusable?
func NewSyncer(
	l *zap.SugaredLogger,
	client *Client,
	opts SyncOptions,
) *Syncer {
	return &Syncer{
		l:    l,
		c:    client,
		opts: opts,

		issuesC: make(chan *github.Issue, opts.DetailsFetchWorkers),
		prC:     make(chan *github.PullRequest, opts.DetailsFetchWorkers),

		outC: make(chan *host.Item, opts.OutputBufferSize),
	}
}

// Sync pulls all issues from its configured repository and pipes them to the
// returned channels. It can only be called once.
func (s *Syncer) Sync(ctx context.Context, wg *sync.WaitGroup) (<-chan *host.Item, <-chan error) {
	// guard against reuse
	if s.used {
		var errC = make(chan error, 1)
		errC <- errors.New("syncer cannot be reused")
		close(errC)
		return nil, errC
	}
	s.used = true

	// execute sync
	s.l.Info("executing sync")
	var errC = s.sync(ctx, wg)
	go func() {
		wg.Wait()
		s.l.Info("sync done, closing output")
		close(s.outC)
	}()
	return s.outC, errC
}

func (s *Syncer) sync(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	wg.Add(2)
	go s.handleIssues(ctx, wg)
	go s.handlePullRequests(ctx, wg)

	// start sync
	var errC = make(chan error)
	wg.Add(2)
	go func() {
		if err := s.c.GetIssues(
			ctx,
			s.opts.Repo.Owner,
			s.opts.Repo.Name,
			s.opts.Filter,
			s.issuesC,
			wg); err != nil {
			errC <- err
		}
	}()
	go func() {
		if err := s.c.GetPullRequests(
			ctx,
			s.opts.Repo.Owner,
			s.opts.Repo.Name,
			s.opts.Filter,
			s.prC,
			wg); err != nil {
			errC <- err
		}
	}()
	go func() {
		s.l.Infow("waiting for waitgroup")
		wg.Wait()
		s.l.Infow("waitgroup done, closing errors")
		close(errC)
	}()

	return errC
}

func (s *Syncer) handleIssues(ctx context.Context, wg *sync.WaitGroup) {
	for i := range s.issuesC {
		var item = &host.Item{
			GitHubID: int(i.GetID()),
			Number:   int(i.GetNumber()),
			Type:     host.ItemTypeIssue,

			Author: i.GetUser().GetName(),
			Opened: i.GetCreatedAt(),
			Closed: i.ClosedAt,

			Title: i.GetTitle(),
			Body:  i.GetBody(),

			// TODO: flesh this out
			Details: map[string]interface{}{
				"comments": i.GetComments(),
			},
		}
		item.WithGitHubReactions(i.Reactions)
		item.WithGitHubLabels(i.Labels)
		s.outC <- item
	}
	s.l.Infow("all issues processed")
	wg.Done()
}

func (s *Syncer) handlePullRequests(ctx context.Context, wg *sync.WaitGroup) {
	for pr := range s.prC {
		var item = &host.Item{
			GitHubID: int(pr.GetID()),
			Number:   int(pr.GetNumber()),
			Type:     host.ItemTypePR,

			Author: pr.GetUser().GetName(),
			Opened: pr.GetCreatedAt(),
			Closed: pr.ClosedAt,

			Title: pr.GetTitle(),
			Body:  pr.GetBody(),

			// TODO: flesh out PR stuff
			Details: map[string]interface{}{
				"commit":      pr.GetMergeCommitSHA(),
				"comments":    pr.GetComments(),
				"commits":     pr.GetCommits(),
				"files":       pr.GetChangedFiles(),
				"additions":   pr.GetAdditions(),
				"deletions":   pr.GetDeletions(),
				"mergability": pr.GetMergeableState(),
			},
		}
		// PR list api does not return reactions :(
		// item.WithGitHubReactions(pr.Reactions)

		// PR list api returns labels as []*github.Label, convert it first
		labels := make([]github.Label, len(pr.Labels))
		for i, l := range pr.Labels {
			if l != nil {
				labels[i] = *l
			}
		}
		item.WithGitHubLabels(labels)
		s.outC <- item
	}
	s.l.Infow("all pull requests processed")
	wg.Done()
}
