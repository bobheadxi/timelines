package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/config"
	"go.bobheadxi.dev/zapx/zpgx"
)

// Database is a low-level wrapper around the database driver
type Database struct {
	l  *zap.SugaredLogger
	pg *pgx.ConnPool
}

// New instantiates a new database
func New(l *zap.SugaredLogger, name string, opts config.Database) (*Database, error) {
	// set up configuration
	var connConfig pgx.ConnConfig
	if opts.PostgresConnURL != "" {
		l.Info("parsing conn string")
		var err error
		connConfig, err = pgx.ParseURI(opts.PostgresConnURL)
		if err != nil {
			return nil, fmt.Errorf("failed to read db conn url: %v", err)
		}
	} else {
		l.Info("using provided parameters")
		port, _ := strconv.Atoi(opts.Port)
		connConfig = pgx.ConnConfig{
			Host:     opts.Host,
			Port:     uint16(port),
			Database: opts.Database,

			// authentication
			User:      opts.User,
			Password:  opts.Password,
			TLSConfig: opts.TLS,

			// misc metadata
			RuntimeParams: map[string]string{
				"application_name": name,
			},
			Logger: zpgx.NewLogger(l.Desugar().Named("px"), zpgx.Options{
				LogInfoAsDebug: true,
			}),
		}
	}
	l.Infow("set up configuration",
		"host", connConfig.Host,
		"database", connConfig.Database)

	// init connection pool
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %v", err)
	}

	// create struct
	var db = &Database{
		pg: pool,
		l:  l,
	}

	// set up statements and whatnot
	db.Repos().init()

	return db, nil
}

// Repos instantaite a new ReposDatabase client
func (db *Database) Repos() *ReposDatabase {
	return &ReposDatabase{db: db, l: db.l.Named("repos")}
}

// Pool returns the underlying pgx connection pool
func (db *Database) Pool() *pgx.ConnPool { return db.pg }

// Close disconnects from the database
func (db *Database) Close() { db.pg.Close() }
