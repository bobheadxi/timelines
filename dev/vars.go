package dev

import "github.com/bobheadxi/projector/config"

var (
	StoreOptions = config.Store{
		Address:  "127.0.0.1:6379",
		Password: "",
	}

	DatabaseOptions = config.Database{
		Host:     "127.0.0.1",
		Port:     "5431",
		Database: "projector-dev",
		User:     "bobheadxi",
	}
)
