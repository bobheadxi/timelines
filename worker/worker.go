package worker

import (
	"github.com/bobheadxi/projector/log"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gocloud.dev/server"
)

// RunOpts denotes worker options
type RunOpts struct {
	Port string
}

// Run spins up the worker
func Run(
	l *zap.SugaredLogger,
	opts RunOpts,
) error {
	// init server with diagnostic hooks
	var srv = server.New(&server.Options{
		// TODO
		RequestLogger: log.NewRequestLogger(l.Named("requests")),
	})

	var mux = chi.NewMux()
	// TODO

	l.Infow("spinning up worker",
		"port", opts.Port)
	return srv.ListenAndServe(":"+opts.Port, mux)
}
