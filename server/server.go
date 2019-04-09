package server

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gocloud.dev/server"

	"github.com/bobheadxi/res"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/store"
)

// RunOpts denotes server options
type RunOpts struct {
	Port     string
	Store    config.Store
	Database config.Database
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

	// init clients
	store, err := store.NewClient(l.Named("store"), "server", opts.Store)
	if err != nil {
		return err
	}
	defer store.Close()
	database, err := db.New(l.Named("db"), "timelines.worker", opts.Database)
	if err != nil {
		return err
	}
	defer database.Close()

	// init handlers
	var resolver = newResolver(l.Named("resolver"))
	var webhook = newWebhookHandler(l.Named("webhooks"), database, store)

	// set up endpoints
	var mux = chi.NewMux()
	mux.Route("/api", func(r chi.Router) {
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			res.R(w, r, res.MsgOK("the timelines api is online!"))
		})
		r.Handle("/playground", handler.Playground("timelines API Playground", "/api/query"))
		r.Handle("/query", handler.GraphQL(timelines.NewExecutableSchema(timelines.Config{
			Resolvers: resolver,
		})))
	})
	mux.Route("/webhooks", func(r chi.Router) {
		r.HandleFunc("/github", webhook.handleGitHub)
	})

	// let's go!
	l.Infow("spinning up server",
		"port", opts.Port)
	return srv.ListenAndServe(":"+opts.Port, mux)
}
