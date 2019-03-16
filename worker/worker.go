package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/analysis"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/git"
	"github.com/bobheadxi/timelines/host"
	"github.com/bobheadxi/timelines/host/github"
	"github.com/bobheadxi/timelines/store"
)

// RunOpts denotes worker options
type RunOpts struct {
	Workers  int
	Store    config.Store
	Database config.Database
}

// Run spins up the worker
func Run(
	l *zap.SugaredLogger,
	stop chan bool,
	opts RunOpts,
) error {
	store, err := store.NewClient(l.Named("store"), opts.Store)
	if err != nil {
		return err
	}
	defer store.Close()

	database, err := db.New(l.Named("db"), "timelines.worker", opts.Database)
	if err != nil {
		return err
	}
	defer database.Close()

	hub, err := github.NewSigningClient(l.Named("github-signer"), github.NewEnvAuth())
	if err != nil {
		return err
	}

	repos := git.NewManager(l.Named("git"), git.ManagerOpts{
		Workdir: "./tmp",
	})

	// set up worker
	var (
		errC = make(chan error, 10)
		w    = newWorker(
			l,
			store,
			database,
			hub,
			repos,
			errC)
	)

	// let's go!
	defer close(errC)
	l.Infow("spinning up worker")
	if opts.Workers == 0 {
		opts.Workers = 1
	}
	for i := 0; i < opts.Workers; i++ {
		go w.processJobs(stop, errC)
	}
	for {
		select {
		case err := <-errC:
			w.l.Errorw("critical error encountered - stopping worker",
				"error", err)
			stop <- true
			return err
		case <-stop:
			w.l.Info("stopping worker")
			return nil
		}
	}
}

type worker struct {
	store *store.Client
	db    *db.Database

	hub *github.SigningClient
	git *git.Manager

	errC chan<- error

	l *zap.SugaredLogger
}

func newWorker(
	l *zap.SugaredLogger,
	store *store.Client,
	db *db.Database,
	hub *github.SigningClient,
	git *git.Manager,
	errC chan<- error,
) *worker {
	return &worker{
		store,
		db,
		hub,
		git,
		errC,
		l,
	}
}

func (w *worker) processJobs(stop <-chan bool, errC chan<- error) {
	for {
		jobC, jobErrC := w.store.RepoJobs().Dequeue(30 * time.Second)
		select {
		case <-stop:
			w.l.Info("stopping job processor")
			return
		case err := <-jobErrC:
			w.l.Errorw("error received when dequeing",
				"error", err)
			continue
		case job := <-jobC:
			if job == nil {
				continue
			}
			w.l.Info("job dequeued", "job.id", job.ID)
			w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
				GitHubSync: &store.StateMeta{
					State: store.StateInProgress,
				},
				Analysis: &store.StateMeta{
					State: store.StateInProgress,
				},
			})

			var (
				wg  sync.WaitGroup
				ctx = context.Background() // TODO: enforce timeout?
			)

			// spin up handlers and wait until completion
			wg.Add(2)
			go w.githubSync(ctx, job, &wg)
			go w.gitAnalysis(ctx, job, &wg)
			wg.Wait()
			w.l.Infow("job processed", "job.id", job.ID)
		}
	}
}

func (w *worker) gitAnalysis(ctx context.Context, job *store.RepoJob, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		l      = w.l.With("job.id", job.ID).Named("git_analysis")
		start  = time.Now()
		remote = fmt.Sprintf("https://github.com/%s/%s.git", job.Owner, job.Repo)
	)

	// load or download repo
	repo, err := w.git.Load(ctx, remote)
	if err != nil {
		repo, err = w.git.Download(ctx, remote, git.DownloadOpts{})
		if err != nil {
			w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
				Analysis: &store.StateMeta{
					State:   store.StateError,
					Message: err.Error(),
				},
			})
			l.Errorw("repo does not exist and could not download", "error", err)
			return
		}
	}

	// set up analysis
	an, err := analysis.NewGitAnalyser(
		l.Named("analyzer"),
		repo.GitRepo(),
		analysis.GitRepoAnalyserOptions{})
	if err != nil {
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: &store.StateMeta{
				State:   store.StateError,
				Message: fmt.Sprintf("analysis.setup: %v", err),
			},
		})
		l.Errorw("analysis failed", "error", err)
		return
	}

	// execute
	report, err := an.Analyze()
	if err != nil {
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: &store.StateMeta{
				State:   store.StateError,
				Message: fmt.Sprintf("analysis.execute: %v", err),
			},
		})
		l.Errorw("analysis failed", "error", err)
		return
	}

	// TODO dump in DB
	l.Infow("analysis complete",
		"duration", time.Since(start),
		"report", report)
	w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
		Analysis: &store.StateMeta{
			State: store.StateDone,
		},
	})
}

func (w *worker) githubSync(ctx context.Context, job *store.RepoJob, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		l     = w.l.With("job.id", job.ID).Named("github_sync")
		start = time.Now()
		repos = w.db.Repos()
	)

	// set up client
	client, err := w.hub.GetInstallationClient(ctx, job.InstallationID)
	if err != nil {
		l.Errorw("failed to authenticate for installation",
			"error", err,
			"github.installation_id", job.InstallationID)
		return
	}

	// check for entry in DB
	repoID, err := repos.GetRepositoryID(ctx, job.Owner, job.Repo)
	if err != nil || repoID == 0 {
		// repo must exist at this point, since server must create it first
		l.Errorw("could not find repository entry in database",
			"error", err,
			"repository", job.Owner+"/"+job.Repo)
		w.errC <- fmt.Errorf("could not find ID for repository '%s/%s'", job.Owner, job.Repo)
		return
	}
	l = l.With("db.id", repoID)

	// init syncer
	var syncer = github.NewSyncer(
		w.l.Named("github_sync").With("job.id", job.ID),
		client,
		github.SyncOptions{
			Repo: github.Repo{
				Owner: job.Owner,
				Name:  job.Repo,
			},
			Filter: github.ItemFilter{
				State:    github.IssueStateAll,
				Interval: time.Second,
			},
			DetailsFetchWorkers: 3,
			OutputBufferSize:    3,
		})

	var (
		syncWG     sync.WaitGroup
		stopPipe   = make(chan bool)
		itemCount  int32
		errorCount int32
		bufsize    = 30

		itemsC, syncErrC = syncer.Sync(ctx, &syncWG)
	)
	go func() {
		var (
			cur int
			buf = make([]*host.Item, bufsize)
		)
		defer func() {
			// if the first item of buffer is non-nil, there are some number of items
			// that needs to be dumped
			if buf[0] != nil {
				if err := repos.InsertGitHubItems(ctx, repoID, buf); err != nil {
					l.Errorw("failed to clear github items", "error", err)
					w.errC <- err
					atomic.AddInt32(&errorCount, 1)
					return
				}
				l.Infow("buffer cleared")
			}
		}()
		for {
			select {
			case item := <-itemsC:
				atomic.AddInt32(&itemCount, 1)
				if item == nil {
					itemsC = nil
					continue
				}

				// if we're at buffer limit, dump buffer
				if cur >= bufsize {
					var err = repos.InsertGitHubItems(ctx, repoID, buf)
					// clear buffer straight away to prevent defer from double-inserting
					cur = 0
					buf = nil
					buf = make([]*host.Item, bufsize)
					if err != nil {
						l.Errorw("failed to insert github items", "error", err)
						w.errC <- err
						atomic.AddInt32(&errorCount, 1)
						return
					}
					l.Infow("buffer cleared")
				}

				// insert item into buffer
				buf[cur] = item
				cur++

			case err := <-syncErrC:
				if err != nil {
					l.Errorw("error occured while syncing",
						"error", err)
					w.errC <- err
					atomic.AddInt32(&errorCount, 1)
				} else {
					syncErrC = nil
				}
				return

			case <-stopPipe:
				return

			case <-ctx.Done():
				return
			}
		}
	}()

	// wait until sync wraps up
	syncWG.Wait()
	l.Infow("github sync complete",
		"items", itemCount,
		"duration", time.Since(start))
	stopPipe <- true

	// update job state
	var state *store.RepoJobState
	if errorCount > 0 {
		state = &store.RepoJobState{
			GitHubSync: &store.StateMeta{
				State:   store.StateError,
				Message: "an error was encountered during sync", // TODO: track actual error
			},
		}
	} else {
		state = &store.RepoJobState{
			GitHubSync: &store.StateMeta{
				State: store.StateDone,
				Meta:  map[string]interface{}{"items": itemCount},
			},
		}
	}
	w.store.RepoJobs().SetState(job.ID, state)
	l.Infow("state updated", "state", state)
}
