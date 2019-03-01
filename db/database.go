package db

import (
	"crypto/tls"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

// Database is a low-level wrapper around the database driver
type Database struct {
	l  *zap.SugaredLogger
	pg *pg.DB
}

// Options denotes database instantiation options
type Options struct {
	Address  string
	Database string

	TLS      *tls.Config
	User     string
	Password string
}

// New instantiates a new database
func New(l *zap.SugaredLogger, opts Options) (*Database, error) {
	var driver = pg.Connect(&pg.Options{
		ApplicationName: "projector",

		Addr:     opts.Address,
		Database: opts.Database,

		TLSConfig: opts.TLS,
		User:      opts.User,
		Password:  opts.Password,
	})
	var db = &Database{
		pg: driver,
		l:  l,
	}
	return db, db.init()
}

// init runs any required initialization on the database
func (db *Database) init() error {
	var now = time.Now()
	for model, opts := range map[interface{}]*orm.CreateTableOptions{
		&Repository{}: &orm.CreateTableOptions{
			FKConstraints: true,
		},
	} {
		if err := db.pg.CreateTable(model, opts); err != nil {
			db.l.Errorw("could not create table",
				"error", err,
				"model", model)
			return err
		}
	}
	db.l.Infow("table instantiated",
		"duration", time.Since(now))
	return nil
}
