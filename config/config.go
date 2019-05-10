package config

import (
	"github.com/BurntSushi/toml"
)

var Conf = &Config{}

type Config struct {
	Net      string
	Port     string
	MySQLUrl string
	Seeds    []string
	NeoUrl   string
	Log      struct {
		OutputWay string
		Output    string
	}
	OpenSync         bool
	OnlySyncBlock    bool
	SyncBlockThreads int
}

func Init(cfgName string) (err error) {
	_, err = toml.DecodeFile(cfgName, &Conf)
	return
}
