package github

import (
	"context"
	"errors"
	"sync"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
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

	issuesC       chan *github.Issue
	fetchDetailsC chan *github.Issue

	outC chan *Item

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

		issuesC:       make(chan *github.Issue, opts.DetailsFetchWorkers),
		fetchDetailsC: make(chan *github.Issue, opts.DetailsFetchWorkers),

		outC: make(chan *Item, opts.OutputBufferSize),
	}
}

// Sync pulls all issues from its configured repository and pipes them to the
// returned channels. It can only be called once.
func (s *Syncer) Sync(ctx context.Context, wg *sync.WaitGroup) (<-chan *Item, <-chan error) {
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
	// spin up workers
	for i := 0; i < s.opts.DetailsFetchWorkers; i++ {
		wg.Add(1)
		go s.fetchDetails(ctx, wg)
	}
	wg.Add(1)
	go s.handleIssues(ctx, wg)

	// start sync
	var errC = make(chan error)
	wg.Add(1)
	go func() {
		if err := s.c.GetIssues(
			ctx,
			s.opts.Repo.Owner,
			s.opts.Repo.Name,
			s.opts.Filter,
			s.issuesC,
			s.fetchDetailsC,
			wg); err != nil {
			errC <- err
		}
		s.l.Infow("all issues loaded done, waiting")
		wg.Wait()
		s.l.Infow("sync done, closing error output")
		close(errC)
	}()

	return errC
}

func (s *Syncer) handleIssues(ctx context.Context, wg *sync.WaitGroup) {
	for i := range s.issuesC {
		var item = &Item{
			GitHubID: int(i.GetID()),
			Number:   int(i.GetNumber()),
			Type:     ItemTypeIssue,

			Author: i.GetUser().GetName(),
			Opened: i.GetCreatedAt(),
			Closed: i.ClosedAt,

			Title: i.GetTitle(),
			Body:  i.GetBody(),

			// TODO: flesh this out
			Details: map[string]interface{}{},
		}
		item.WithReactions(i.Reactions)
		item.WithLabels(i.Labels)
		s.outC <- item
	}
	s.l.Infow("all issues processed")
	wg.Done()
}

func (s *Syncer) fetchDetails(ctx context.Context, wg *sync.WaitGroup) {
	for i := range s.fetchDetailsC {
		if pr, err := s.c.GetPullRequest(ctx, i); err != nil {
			s.l.Errorw("failed to get pull request",
				"issue", i.GetNumber())
		} else {
			var item = &Item{
				GitHubID: int(pr.GetID()),
				Number:   int(i.GetNumber()),
				Type:     ItemTypePR,

				Author: i.GetUser().GetName(),
				Opened: i.GetCreatedAt(),
				Closed: i.CreatedAt,

				Title: i.GetTitle(),
				Body:  i.GetBody(),

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
			item.WithReactions(i.Reactions)
			item.WithLabels(i.Labels)
			s.outC <- item
		}
	}
	s.l.Infow("all detail fetching processed")
	wg.Done()
}
