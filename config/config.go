package config

import (
	"github.com/BurntSushi/toml"
)

var Conf = &Config{}

type Config struct {
	Net      string
	MySQLUrl string
	Seeds    []string
	NeoUrl   string
}

func Init(cfgName string) (err error) {
	_, err = toml.DecodeFile(cfgName, &Conf)
	return
}
