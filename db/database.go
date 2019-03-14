package db

import (
	"strconv"

	// need https://github.com/jackc/pgx/issues/335
	"github.com/bobheadxi/pgx"

	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/log"
)

// Database is a low-level wrapper around the database driver
type Database struct {
	l  *zap.SugaredLogger
	pg *pgx.ConnPool
}

// New instantiates a new database
func New(l *zap.SugaredLogger, name string, opts config.Database) (*Database, error) {
	port, _ := strconv.Atoi(opts.Port)
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgx.ConnConfig{
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

		// TODO
		Logger: log.NewDatabaseLogger(l),
	}})
	if err != nil {
		return nil, err
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
