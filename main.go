package main

import (
	"./config"
	"github.com/sirupsen/logrus"
)

func main(){

	logrus.Info("Initializing Git Sync...")

	for _, repo := range config.Config.Repositories {
		logrus.Infof("Detected repo %s with %d mirrors", repo.Name, len(repo.Mirrors));
	}

}
