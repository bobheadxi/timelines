package server

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen-contrib/gqlapollotracing"
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/bobheadxi/res"
	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/graphql/go/timelines"
	"github.com/bobheadxi/timelines/log"
	"github.com/bobheadxi/timelines/store"
)

// TODO: when to enable, when to disable, do we even want this?
func gqlTracers() handler.Option {
	return func(c *handler.Config) {
		handler.RequestMiddleware(gqlapollotracing.RequestMiddleware())(c)
		handler.Tracer(gqlapollotracing.NewTracer())(c)
	}
}

// RunOpts denotes server options
type RunOpts struct {
	Port     string
	Store    config.Store
	Database config.Database

	Build string
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
		resolver = newRootResolver(l.Named("resolver"), database, store)
		webhook  = newWebhookHandler(l.Named("webhooks"), database, store)
		mux      = chi.NewMux()
		srv      = http.Server{
			Addr:    ":" + opts.Port,
			Handler: mux,
		}
	)

	// set up endpoints
	mux.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer)
	mux.Route("/query", func(r chi.Router) {
		// TODO: improve configuration
		r.Use(
			cors.New(cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			}).Handler,
			log.GraphCtxHandler)
		r.Handle("/", handler.GraphQL(
			timelines.NewExecutableSchema(timelines.Config{
				Resolvers:  resolver,
				Directives: timelines.DirectiveRoot{},
			}),
			gqlTracers(),
			handler.RequestMiddleware(log.NewGraphLogger(l.Desugar().Named("graph"))),
		))
	})
	mux.Route("/webhooks", func(r chi.Router) {
		r.Use(log.NewHTTPLogger(l.Named("webhooks")))
		r.HandleFunc("/github", webhook.handleGitHub)
	})
	mux.Route("/playground", func(r chi.Router) {
		r.Use(log.NewHTTPLogger(l.Named("playground")))
		r.Handle("/", handler.Playground("timelines API Playground", "/query"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res.R(w, r, res.MsgOK("API server is online!",
			"build", opts.Build))
	})

	// let's go!
	l.Infow("spinning up server",
		"port", opts.Port)
	go func() {
		<-stop
		l.Info("shutting down server")
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// rootResolver implements the timelines GraphQL API
type rootResolver struct {
	l *zap.SugaredLogger

	q timelines.QueryResolver
	a timelines.RepositoryAnalyticsResolver
}

func newRootResolver(
	l *zap.SugaredLogger,
	d *db.Database,
	s *store.Client,
) timelines.ResolverRoot {
	return &rootResolver{
		l: l,
		q: newQueryResolver(l, d),
		a: newAnalyticsResolver(l, d, s),
	}
}

func (r *rootResolver) Query() timelines.QueryResolver                             { return r.q }
func (r *rootResolver) RepositoryAnalytics() timelines.RepositoryAnalyticsResolver { return r.a }
