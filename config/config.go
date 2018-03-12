package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

var Config *Configuration

func init() {

	tomlDataRaw, err := ioutil.ReadFile("config.toml")

	if err != nil {
		logrus.Panic("Failed to read file config.toml ", err)
	}

	tomlData := string(tomlDataRaw)

	var config Configuration
	if _, err := toml.Decode(tomlData, &config); err != nil {
		logrus.Panic("Failed to parse config.toml ", err)
	}

	Config = &config

}
