package db

import "go.uber.org/zap"

// Repository represents a stored repository. TODO
type Repository struct{}

// ReposDatabase is a client for accessing repository-related databases
type ReposDatabase struct {
	db *Database
	l  *zap.SugaredLogger
}

// NewReposDatabase instantaite a new ReposDatabase client
func NewReposDatabase(l *zap.SugaredLogger, d *Database) *ReposDatabase {
	return &ReposDatabase{db: d, l: l}
}
