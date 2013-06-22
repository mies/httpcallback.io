package main

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Host  HostInformation
	Mongo MongoInformation
}

func OpenConfig(path string) (*Configuration, error) {
	var cfg Configuration
	_, err := toml.DecodeFile(path, &cfg)

	return &cfg, err
}

type HostInformation struct {
	Hostname string
}

type MongoInformation struct {
	ServerUrl    string
	DatabaseName string
}
