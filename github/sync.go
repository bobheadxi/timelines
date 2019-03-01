package github

import (
	"context"
	"strconv"
	"sync"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
)

// Item is an item due for indexing
type Item struct {
	ID   string
	Type string
	Data interface{}
}

// SyncOptions denotes options for a syncer
type SyncOptions struct {
	Repo                Repo
	Filter              ItemFilter
	DetailsFetchWorkers int
	IndexC              chan *Item
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
}

// NewSyncer instantiates a new GitHub Syncer
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
	}
}

// Sync pulls all issues from its configured repository and blocks until done
func (s *Syncer) Sync(ctx context.Context, wg *sync.WaitGroup) error {
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
	}()

	cancellableCtx, cancel := context.WithCancel(ctx)
	for i := 0; i < s.opts.DetailsFetchWorkers; i++ {
		go s.fetchDetails(cancellableCtx, wg)
	}

	wg.Wait()
	cancel()
	select {
	case err := <-errC:
		return err
	default:
		return nil
	}
}

func (s *Syncer) fetchDetails(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case i := <-s.issuesC:
			wg.Add(1)
			s.opts.IndexC <- &Item{
				ID:   strconv.Itoa(int(i.GetID())),
				Type: "issue",
				Data: i,
			}
			wg.Done()
		case i := <-s.fetchDetailsC:
			wg.Add(1)
			if pr, err := s.c.GetPullRequest(ctx, i); err != nil {
				s.l.Errorw("failed to get pull request",
					"issue", i.GetNumber())
			} else {
				s.opts.IndexC <- &Item{
					ID:   strconv.Itoa(int(pr.GetID())),
					Type: "pull-request",
					Data: pr,
				}
			}
			wg.Done()
		}
	}
}
