package worker

import (
	"os"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/git"
	"github.com/bobheadxi/timelines/host/github"
	"github.com/bobheadxi/timelines/store"
)

// RunOpts denotes worker options
type RunOpts struct {
	Name     string
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
	if opts.Name == "" {
		opts.Name = os.Getenv("HOSTNAME")
	}

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
		// this channel is for critical errors
		errC = make(chan error, 10)
		w    = newWorker(
			opts.Name,
			l,
			store,
			database,
			hub,
			repos,
			errC)
	)

	// let's go!
	defer close(errC)
	if opts.Workers == 0 {
		opts.Workers = 1
	}
	l.Infow("spinning up worker processes",
		"workers", opts.Workers)
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
