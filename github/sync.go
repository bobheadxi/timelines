package github

import (
	"context"
	"errors"
	"sync"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
)

// ItemType denotes supported GitHub item types
type ItemType string

const (
	// ItemTypeIssue is a GitHub issue
	ItemTypeIssue ItemType = "issue"
	// ItemTypePR is a GitHub pull request
	ItemTypePR ItemType = "pull-request"
)

// Item is a GitHub item due for indexing
// TODO: this needs to be better
type Item struct {
	ID     int
	Number int
	Type   ItemType
	Data   interface{}
}

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

// Sync pulls all issues from its configured repository and blocks until done.
// It can only be called once.
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
	var errC = s.sync(ctx, wg)
	go func() {
		wg.Wait()
		close(s.outC)
	}()
	return s.outC, errC
}

func (s *Syncer) sync(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	// spin up workers
	for i := 0; i < s.opts.DetailsFetchWorkers; i++ {
		go s.fetchDetails(ctx, wg)
	}
	go s.handleIssues(ctx, wg)

	// start sync
	var errC = make(chan error)
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
		wg.Wait()
		close(errC)
	}()

	return errC
}

func (s *Syncer) handleIssues(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for i := range s.issuesC {
		s.outC <- &Item{
			ID:     int(i.GetID()),
			Number: int(i.GetNumber()),
			Type:   ItemTypeIssue,
			Data:   i,
		}
	}
	wg.Done()
}

func (s *Syncer) fetchDetails(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for i := range s.fetchDetailsC {
		if pr, err := s.c.GetPullRequest(ctx, i); err != nil {
			s.l.Errorw("failed to get pull request",
				"issue", i.GetNumber())
		} else {
			s.outC <- &Item{
				ID:     int(pr.GetID()),
				Number: int(pr.GetNumber()),
				Type:   ItemTypePR,
				Data:   pr,
			}
		}
	}
	wg.Done()
}
