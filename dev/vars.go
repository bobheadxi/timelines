package dev

import "github.com/bobheadxi/projector/config"

var (
	StoreOptions = config.Store{
		Address:  "127.0.0.1:6379",
		Password: "",
	}

	DatabaseOptions = config.Database{
		Address:  "127.0.0.1:5431",
		Database: "projector_dev",
		User:     "bobheadxi",
		Password: "bobheadxi",
	}
)
