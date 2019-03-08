package db

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"

	"github.com/bobheadxi/projector/config"
)

// Database is a low-level wrapper around the database driver
type Database struct {
	l  *zap.SugaredLogger
	pg *pg.DB
}

// New instantiates a new database
func New(l *zap.SugaredLogger, opts config.Database) (*Database, error) {
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

// Repos instantaite a new ReposDatabase client
func (db *Database) Repos() *ReposDatabase {
	return &ReposDatabase{db: db, l: db.l.Named("repos")}
}

// init runs any required initialization on the database
func (db *Database) init() error {
	var now = time.Now()
	for model, opts := range map[interface{}]*orm.CreateTableOptions{
		&Repository{}: &orm.CreateTableOptions{
			FKConstraints: true,
			IfNotExists:   true,
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
