package server

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/bobheadxi/res"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
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
	var (
		resolver = newRootResolver(l.Named("resolver"), database)
		webhook  = newWebhookHandler(l.Named("webhooks"), database, store)
		mux      = chi.NewMux()
		srv      = http.Server{
			Addr:    ":" + opts.Port,
			Handler: mux,
		}
	)

	// set up endpoints
	mux.Handle("/playground", handler.Playground("timelines API Playground", "/query"))
	mux.Route("/query", func(r chi.Router) {
		// TODO: improve configuration
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}).Handler)
		r.Handle("/", handler.GraphQL(timelines.NewExecutableSchema(timelines.Config{
			Resolvers: resolver,
		})))
	})
	mux.Route("/webhooks", func(r chi.Router) {
		r.HandleFunc("/github", webhook.handleGitHub)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res.R(w, r, res.MsgOK("API server is online!"))
	})

	// let's go!
	l.Infow("spinning up server",
		"port", opts.Port)
	go func() {
		<-stop
		l.Info("shutting down server")
		srv.Shutdown(context.Background())
	}()
	return srv.ListenAndServe()
}

// rootResolver implements the timelines GraphQL API
type rootResolver struct {
	l *zap.SugaredLogger
	q timelines.QueryResolver
}

func newRootResolver(
	l *zap.SugaredLogger,
	database *db.Database,
) *rootResolver {
	return &rootResolver{
		l: l,
		q: &queryResolver{
			db: database,
			l:  l.Named("query"),
		},
	}
}

func (r *rootResolver) Query() timelines.QueryResolver { return r.q }
