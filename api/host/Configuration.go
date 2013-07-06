package host

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Host   *HostInformation
	Mongo  *MongoInformation
	Github *OAuthProviderInfo
}

func OpenConfig(path string) (*Configuration, error) {
	var cfg Configuration
	_, err := toml.DecodeFile(path, &cfg)

	return &cfg, err
}

type OAuthProviderInfo struct {
	ClientId       string
	ClientSecret   string
	AuthorizeUrl   string
	AccessTokenUrl string
}

type HostInformation struct {
	Hostname string
}

type MongoInformation struct {
	UseMongo     bool
	ServerUrl    string
	DatabaseName string
}
