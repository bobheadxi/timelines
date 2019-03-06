package server

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gocloud.dev/server"

	"github.com/bobheadxi/projector/graphql/golang/api"
)

// RunOpts denotes server options
type RunOpts struct {
	Port string
}

// Run spins up the server
func Run(l *zap.SugaredLogger, resolver *Resolver, opts RunOpts) error {
	// init server with diagnostic hooks
	var srv = server.New(&server.Options{
		// TODO
		// RequestLogger: l.Named("requests"),
	})

	// set up endpoints
	var mux = chi.NewMux()
	mux.Handle("/", handler.Playground("Projector API Playground", "/query"))
	mux.Handle("/query", handler.GraphQL(api.NewExecutableSchema(api.Config{
		Resolvers: resolver,
	})))

	// let's go!
	l.Infow("spinning up server",
		"port", opts.Port)
	return srv.ListenAndServe(":"+opts.Port, mux)
}
