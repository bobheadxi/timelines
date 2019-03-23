package config

import (
	"crypto/tls"
	"os"
)

// Store denotes store client instantiation options
type Store struct {
	Address  string
	Password string

	TLS *tls.Config
}

// NewStoreConfig loads store configuration from environment
func NewStoreConfig() Store {
	return Store{
		Address:  os.Getenv(""),
		Password: os.Getenv(""),

		TLS: nil,
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
}

// NewDatabaseConfig loads database configuration from environment
func NewDatabaseConfig() Database {
	return Database{
		Host:     os.Getenv(""),
		Port:     os.Getenv(""),
		Database: os.Getenv(""),

		User:     os.Getenv(""),
		Password: os.Getenv(""),

		Drop: false,
		TLS:  nil,
	}
}
