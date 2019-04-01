package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/analysis"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/git"
	"github.com/bobheadxi/timelines/host"
	"github.com/bobheadxi/timelines/host/github"
	"github.com/bobheadxi/timelines/store"
)

type worker struct {
	name string

	store *store.Client
	db    *db.Database

	hub *github.SigningClient
	git *git.Manager

	errC chan<- error

	l *zap.SugaredLogger
}

func newWorker(
	name string,
	l *zap.SugaredLogger,
	store *store.Client,
	db *db.Database,
	hub *github.SigningClient,
	git *git.Manager,
	errC chan<- error,
) *worker {
	return &worker{
		name,
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
		l = w.l.
			With("job.id", job.ID, "job.repo", job.Owner+"/"+job.Repo).
			Named("git_analysis")
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
					Message: fmt.Sprintf("analysis.git_clone: %v", err),
					Meta:    map[string]interface{}{"worker": w.name},
				},
			})
			l.Errorw("repo does not exist and could not download", "error", err)
			return
		}
	}
	l.Infow("repo loaded", "dir", repo.Dir())

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
				Meta:    map[string]interface{}{"worker": w.name},
			},
		})
		l.Errorw("analysis failed", "error", err)
		return
	}

	// execute
	l.Info("executing analysis")
	report, err := an.Analyze()
	if err != nil {
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: &store.StateMeta{
				State:   store.StateError,
				Message: fmt.Sprintf("analysis.execute: %v", err),
				Meta: map[string]interface{}{
					"duration": time.Since(start),
					"worker":   w.name,
				},
			},
		})
		l.Errorw("analysis failed", "error", err)
		return
	}

	dur := time.Since(start)
	l.Infow("analysis complete",
		"duration", dur,
		"report", report) // TODO dump in DB instead
	w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
		Analysis: &store.StateMeta{
			State: store.StateDone,
			Meta: map[string]interface{}{
				"duration": dur,
				"worker":   w.name,
			},
		},
	})
}

func (w *worker) githubSync(ctx context.Context, job *store.RepoJob, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		l = w.l.
			With("job.id", job.ID, "job.repo", job.Owner+"/"+job.Repo).
			Named("github_sync")
		start = time.Now()
		repos = w.db.Repos()
	)

	// set up client
	client, err := w.hub.GetInstallationClient(ctx, job.InstallationID)
	if err != nil {
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: &store.StateMeta{
				State:   store.StateError,
				Message: fmt.Sprintf("github_sync.new_client: %v", err),
				Meta:    map[string]interface{}{"worker": w.name, "installation": job.InstallationID},
			},
		})
		l.Errorw("failed to authenticate for installation",
			"error", err,
			"github.installation_id", job.InstallationID)
		return
	}

	// check for entry in DB
	repoID, err := repos.GetRepositoryID(ctx, job.Owner, job.Repo)
	if err != nil || repoID == 0 {
		// repo must exist at this point, since server must create it first
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: &store.StateMeta{
				State:   store.StateError,
				Message: fmt.Sprintf("github_sync.get_repo: %v", err),
				Meta:    map[string]interface{}{"worker": w.name},
			},
		})
		l.Errorw("could not find repository entry in database",
			"error", err,
			"repository", job.Owner+"/"+job.Repo)
		return
	}
	l = l.With("db.id", repoID)
	l.Info("repository found")

	// init syncer
	var syncer = github.NewSyncer(
		l.Named("syncer"),
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
	l.Info("initializing sync")
	go func() {
		var (
			cur int
			buf = make([]*host.Item, bufsize)
		)
		defer func() {
			// if the first item of buffer is non-nil, there are some number of items
			// that needs to be dumped
			if buf[0] != nil {
				if err := repos.InsertHostItems(ctx, repoID, buf); err != nil {
					l.Errorw("failed to clear github items", "error", err)
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
					var err = repos.InsertHostItems(ctx, repoID, buf)
					// clear buffer straight away to prevent defer from double-inserting
					cur = 0
					buf = nil
					buf = make([]*host.Item, bufsize)
					if err != nil {
						l.Errorw("failed to insert github items", "error", err)
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
	dur := time.Since(start)
	syncWG.Wait()
	l.Infow("github sync complete",
		"items", itemCount,
		"duration", dur)
	stopPipe <- true

	// update job state
	var state *store.RepoJobState
	meta := map[string]interface{}{
		"items":    itemCount,
		"duration": dur,
		"worker":   w.name,
	}
	if errorCount > 0 {
		state = &store.RepoJobState{
			GitHubSync: &store.StateMeta{
				State:   store.StateError,
				Message: "an error was encountered during sync", // TODO: track actual error
				Meta:    meta,
			},
		}
	} else {
		state = &store.RepoJobState{
			GitHubSync: &store.StateMeta{
				State: store.StateDone,
				Meta:  meta,
			},
		}
	}
	w.store.RepoJobs().SetState(job.ID, state)
	l.Infow("state updated", "state", state)
}
