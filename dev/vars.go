package dev

import (
	"os"

	"github.com/bobheadxi/timelines/config"
)

var (
	StoreOptions = config.Store{
		Address:  "127.0.0.1:6379",
		Password: "",
	}

	DatabaseOptions = config.Database{
		Host:     "127.0.0.1",
		Port:     "5431",
		Database: "timelines-dev",
		User:     "bobheadxi",
	}
)

// GetTestInstallationID returns $GITHUB_TEST_INSTALLTION
func GetTestInstallationID() string { return os.Getenv("GITHUB_TEST_INSTALLTION") }
