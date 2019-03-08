package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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
			stop <- true
			return err
		case <-stop:
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
			return
		case err := <-jobErrC:
			w.l.Errorw("error received when dequeing",
				"error", err)
			continue
		case job := <-jobC:
			if job == nil {
				continue
			}
			var l = w.l.With("job.id", job.ID)
			l.Info("job dequeued")

			client, err := w.hub.GetInstallationClient(context.Background(), job.InstallationID)
			if err != nil {
				l.Errorw("failed to authenticate for installation",
					"error", err,
					"github.installation_id", job.InstallationID)
				continue
			}
			var syncer = github.NewSyncer(
				w.l.Named("github.sync").With("job.id", job.ID),
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

			var wg sync.WaitGroup
			var stopPipe = make(chan bool)
			itemsC, syncErrC := syncer.Sync(context.Background(), &wg)
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
					}
				}
			}()
			wg.Wait()
			l.Infow("job processed",
				"items", count)
			stopPipe <- true
		}
	}
}
