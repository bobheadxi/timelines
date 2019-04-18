package config

import (
	"crypto/tls"
	"os"
)

var (
	// Commit is the commit hash of this build
	Commit string
)

// BuildMeta denotes build metadata
type BuildMeta struct {
	Commit string
}

// NewBuildMeta instantiates a new build metadata struct from the environment.
// Currently leverages Heroku's Dyno Metadata: https://devcenter.heroku.com/articles/dyno-metadata
func NewBuildMeta() BuildMeta {
	return BuildMeta{
		Commit: firstOf(Commit, os.Getenv("HEROKU_SLUG_COMMIT")),
	}
}

// Store denotes store client instantiation options
type Store struct {
	Address  string
	Password string

	TLS *tls.Config

	// this has priority
	RedisConnURL string
}

// NewStoreConfig loads store configuration from environment
func NewStoreConfig() Store {
	return Store{
		Address:  os.Getenv("STORE_ADDRESS"),
		Password: os.Getenv("STORE_PW"),

		TLS: nil,

		RedisConnURL: os.Getenv("REDIS_URL"),
	}
}

// Database denotes database instantiation options
type Database struct {
	Host     string
	Port     string
	Database string

	User     string
	Password string

	Drop bool
	TLS  *tls.Config

	// this has priority
	PostgresConnURL string
}

// NewDatabaseConfig loads database configuration from environment
func NewDatabaseConfig() Database {
	return Database{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),

		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PW"),

		Drop: false,
		TLS:  nil,

		PostgresConnURL: os.Getenv("DATABASE_URL"),
	}
}

func firstOf(vars ...string) string {
	for _, s := range vars {
		if s != "" {
			return s
		}
	}
	return ""
}
