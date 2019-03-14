package server

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gocloud.dev/server"

	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/log"
)

// RunOpts denotes server options
type RunOpts struct {
	Port string
}

// Run spins up the server
func Run(
	l *zap.SugaredLogger,
	stop chan bool,
	opts RunOpts,
) error {
	// init server with diagnostic hooks
	var srv = server.New(&server.Options{
		// TODO
		RequestLogger: log.NewRequestLogger(l.Named("requests")),
	})

	// init resolver
	var res = newResolver()

	// set up endpoints
	var mux = chi.NewMux()
	mux.Route("/api", func(r chi.Router) {
		r.Handle("/", handler.Playground("timelines API Playground", "/api/query"))
		r.Handle("/query", handler.GraphQL(timelines.NewExecutableSchema(timelines.Config{
			Resolvers: res,
		})))
	})

	// let's go!
	l.Infow("spinning up server",
		"port", opts.Port)
	return srv.ListenAndServe(":"+opts.Port, mux)
}
