package dev

import (
	"os"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/host"
)

var (
	// StoreOptions denotes store configuration for use with devenv
	StoreOptions = config.Store{
		Address:  "127.0.0.1:6379",
		Password: "",
	}

	// DatabaseOptions denotes database configuration for use with devenv
	DatabaseOptions = config.Database{
		Host:     "127.0.0.1",
		Port:     "5431",
		Database: "timelines-dev",
		User:     "bobheadxi",
	}

	// SampleRepo is a demo repo
	SampleRepo = &host.BaseRepo{
		Host:  host.HostGitHub,
		Owner: "bobheadxi",
		Name:  "calories",
	}
)

// GetTestInstallationID returns $GITHUB_TEST_INSTALLTION
func GetTestInstallationID() string { return os.Getenv("GITHUB_TEST_INSTALLTION") }
