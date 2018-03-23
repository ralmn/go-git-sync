package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"math/rand"
	"bytes"
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


	for i, repo := range config.Repositories {
		if repo.Token == "" {
			repo.Token = randSeq(10)
		}
		logrus.Infof("Generate auth token '%s' for repository : %s", repo.Token, repo.Name)
		config.Repositories[i] = repo
 	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err == nil {
		ioutil.WriteFile("config.toml", buf.Bytes(), 0700)
	}

	Config = &config

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}