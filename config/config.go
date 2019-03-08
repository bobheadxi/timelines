package config

import "crypto/tls"

// Store denotes store client instantiation options
type Store struct {
	Address string

	TLS      *tls.Config
	Password string
}

// Database denotes database instantiation options
type Database struct {
	Address  string
	Database string

	TLS      *tls.Config
	User     string
	Password string
}
