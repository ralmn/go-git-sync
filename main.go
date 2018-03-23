package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ralmn/go-git-sync/config"
	"github.com/ralmn/go-git-sync/repositories"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
)

var repos map[string]repositories.Repository

var version string

func main() {


	logrus.Infof("Initializing Git Sync (version : %s)...", version)

	repos = map[string]repositories.Repository{}

	for _, repo := range config.Config.Repositories {
		logrus.Infof("Detected repo %s with %d mirrors", repo.Name, len(repo.Mirrors))

		repository := repositories.Repository{Name: repo.Name, Mirrors: repo.Mirrors, Token:repo.Token}
		repository.TryToFirstClone()
		repository.SetupRemotes()
		repos[repo.Name] = repository
	}

	InitWeb()

}



func InitWeb() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/webhook/push/{repoName}/{remoteName}", webhookPush)
	logrus.Fatal(http.ListenAndServe(":8080", router))
}

func webhookPush(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	secretToken := req.URL.Query().Get("secret")

	repositoryName := vars["repoName"]
	remoteName := vars["remoteName"]

	repo, ok := repos[repositoryName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Repository '%s' not found", repositoryName)
		return
	}

	if secretToken != repo.Token {
		res.WriteHeader(403)
		fmt.Fprintf(res,"Secret token no valid")
		return
	}

	err := repo.FetchRemote(remoteName)
	if err != nil {
		if err != git.NoErrAlreadyUpToDate {
			res.WriteHeader(500)
			fmt.Fprint(res, "Failed to fetch : ", err)
			return
		}
	}

	err = repo.PushToAllRemote()

	if err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, "Failed when pushing : ", err)
		return
	}

	res.WriteHeader(200)
	fmt.Fprintf(res, "OK: Finish sync process for repository %s from remote %s", repositoryName, remoteName)

}
