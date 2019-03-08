package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bobheadxi/projector/analysis"

	"go.uber.org/zap"

	"github.com/bobheadxi/projector/config"
	"github.com/bobheadxi/projector/db"
	"github.com/bobheadxi/projector/git"
	"github.com/bobheadxi/projector/github"
	"github.com/bobheadxi/projector/store"
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

	database, err := db.New(l.Named("db"), opts.Database)
	if err != nil {
		return err
	}

	hub, err := github.NewSigningClient(l.Named("github-signer"), github.NewEnvAuth())
	if err != nil {
		return err
	}

	repos := git.NewManager(l.Named("git"), git.ManagerOpts{
		Workdir: "./tmp",
	})

	var w = newWorker(
		l,
		store,
		database,
		hub,
		repos)
	l.Infow("spinning up worker")
	var errC = make(chan error)
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

	l *zap.SugaredLogger
}

func newWorker(
	l *zap.SugaredLogger,
	store *store.Client,
	db *db.Database,
	hub *github.SigningClient,
	git *git.Manager,
) *worker {
	return &worker{
		store,
		db,
		hub,
		git,
		l,
	}
}

func (w *worker) processJobs(stop <-chan bool, errC chan<- error) {
	for {
		jobC, jobErrC := w.store.RepoJobs().Dequeue()
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
				GitHubSync: store.StateInProgress,
				Analysis:   store.StateInProgress,
			})

			// spin up handlers and wait until completion
			var wg sync.WaitGroup
			var ctx = context.Background() // TODO: enforce timeout?
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
				Analysis: store.StateError,
			})
			l.Errorw("repo does not exist and could not download", "error", err)
			return
		}
	}

	// set up and execute analysis
	an := analysis.NewGitAnalyser(
		l.Named("analyzer"),
		repo.GitRepo(),
		analysis.GitRepoAnalyserOptions{})
	report, err := an.Analyze()
	if err != nil {
		w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
			Analysis: store.StateError,
		})
		l.Errorw("analysis failed", "error", err)
		return
	}

	// TODO dump in DB
	l.Infow("analysis complete",
		"duration", time.Since(start),
		"report", report)
	w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
		Analysis: store.StateDone,
	})
}

func (w *worker) githubSync(ctx context.Context, job *store.RepoJob, wg *sync.WaitGroup) {
	defer wg.Done()

	var l = w.l.With("job.id", job.ID).Named("github_sync")
	var start = time.Now()

	client, err := w.hub.GetInstallationClient(context.Background(), job.InstallationID)
	if err != nil {
		l.Errorw("failed to authenticate for installation",
			"error", err,
			"github.installation_id", job.InstallationID)
		return
	}

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

	var syncWG sync.WaitGroup
	var stopPipe = make(chan bool)
	itemsC, syncErrC := syncer.Sync(context.Background(), &syncWG)
	var count int32
	go func() {
		for {
			select {
			case item := <-itemsC:
				atomic.AddInt32(&count, 1)
				l.Infow("item received",
					"item", item)
				// TODO: dump in database

			case err := <-syncErrC:
				if err != nil {
					l.Errorw("error occured while syncing",
						"error", err)
				}

			case <-stopPipe:
				return

			case <-ctx.Done():
				return
			}
		}
	}()
	syncWG.Wait()
	l.Infow("github sync complete",
		"items", count,
		"duration", time.Since(start))
	stopPipe <- true
	w.store.RepoJobs().SetState(job.ID, &store.RepoJobState{
		GitHubSync: store.StateDone,
	})
	l.Infow("state updated")
}
