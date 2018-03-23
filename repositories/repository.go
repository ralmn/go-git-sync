package repositories

import (
	"github.com/ralmn/go-git-sync/config"
	"os"
	"path"
	"gopkg.in/src-d/go-git.v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"path/filepath"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"errors"
	"fmt"
	"strings"
)

type Repository struct {
	Name          string
	Token         string
	BareDirectory string
	Mirrors       []config.Mirror
	mirrorsMaps map[string]config.Mirror
}

func (repo *Repository) TryToFirstClone() {

	if len(repo.Mirrors) == 0 {
		logrus.Error("Canno't clone a repository with 0 mirrors ! ", repo);
		return
	}

	firstMirror := repo.Mirrors[0]
	if repo.BareDirectory == "" {
		repo.BareDirectory = path.Base(firstMirror.Remote)
	}

	if _, err := os.Stat(repo.bareFullPath()); os.IsNotExist(err) {

		authMethod := createAuthMethod(firstMirror)

		if authMethod == nil {
			return
		}

		logrus.Info("Clonning ", firstMirror.Remote, " to ", repo.bareFullPath())
		git.PlainClone(repo.bareFullPath(), true, &git.CloneOptions{
			URL:          firstMirror.Remote,
			RemoteName:   firstMirror.Name,
			SingleBranch: false,
			Auth:         authMethod,

		})

	}
}
func (repo *Repository) SetupRemotes() {

	if repo.mirrorsMaps == nil {
		repo.mirrorsMaps = map[string]config.Mirror{}
	}

	gitRepo, err := repo.readGitRepository()
	if err != nil {
		logrus.Error("Failed to open git repository : ", repo.BareDirectory, " ", err)
		return
	}

	for _, mirror := range repo.Mirrors {
		if remote, err := gitRepo.Remote(mirror.Name); err != nil {
			repo.createRemote(gitRepo, mirror)
		}else{
			if remote.Config().URLs[0] != mirror.Remote {
				logrus.Infof("Update remote %s with new url %s", mirror.Name, mirror.Remote)
				gitRepo.DeleteRemote(mirror.Name)
				repo.createRemote(gitRepo, mirror)
			}
		}
		repo.mirrorsMaps[mirror.Name] = mirror
		//repo.FetchRemote(mirror.Name)
	}
}

func (repo *Repository) createRemote(gitRepo *git.Repository, mirror config.Mirror){
	if _, err := gitRepo.CreateRemote(&gitConfig.RemoteConfig{
		Name: mirror.Name,
		URLs: []string{mirror.Remote},
	}); err != nil {
		logrus.Error("Failed to create remote ", mirror.Name, " ", err)
	} else {
		logrus.Info("New remote added ", mirror.Name)
	}
}

func (repo *Repository) refSpecs(mirror config.Mirror) ([]gitConfig.RefSpec){
	return []gitConfig.RefSpec{
		gitConfig.RefSpec(fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", mirror.Name)),
		gitConfig.RefSpec("+refs/heads/*:refs/heads/*")}
}

func (repo *Repository) readGitRepository() (*git.Repository, error) {
	return git.PlainOpen(repo.bareFullPath())
}

func (repo *Repository) bareFullPath() (string) {
	return path.Join(baresPath(), repo.BareDirectory)
}
func (repo *Repository) FetchRemote(remoteName string) (error) {

	gitRepo, err := repo.readGitRepository()
	if err != nil {
		return errors.New(fmt.Sprint("Failed to read git repository ", repo.Name, " ", err))
	}

	remote, err := gitRepo.Remote(remoteName)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to find remote ", remoteName, " in repository ", repo.Name, " ", err))
	}

	mirror := repo.Mirror(remoteName)
	if mirror == nil {
		return errors.New(fmt.Sprintf("No mirror found for remote %s", remoteName))
	}
	authMethod := createAuthMethod(*mirror)

	return remote.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Auth: authMethod,
		RefSpecs: repo.refSpecs(*mirror),
	})
}
func (repo Repository) PushToAllRemote() (error) {
	gitRepo, err := repo.readGitRepository()
	if err != nil {
		return errors.New(fmt.Sprint("Failed to read git repository ", repo.Name, " ", err))
	}

	//remotes, err := gitRepo.Remotes()
	//if err != nil {
	//	return errors.New(fmt.Sprint("Failed to list remote in git repository ", repo.Name, " ", err))
	//}

	var strErrors []string

	for _, mirror := range repo.Mirrors {

		remoteName := mirror.Name
		authMethod := createAuthMethod(mirror)

		remote, err := gitRepo.Remote(remoteName)
		if err != nil {
			strErrors = append(strErrors, fmt.Sprintf("Failed to find remote %s  : %s", remoteName, err))
			continue
		}

		logrus.Infof("Pushing to remote %s (%s)", mirror.Name, repo.Name)



		err = remote.Push(&git.PushOptions{
			RemoteName: mirror.Name,
			Auth: authMethod,
			RefSpecs: repo.refSpecs(mirror),
		})
		if err != nil {
			if err != git.NoErrAlreadyUpToDate {
				strErrors = append(strErrors, fmt.Sprint("Failed to push in ", remoteName, " in git repository ", repo.Name, " ", err))
			}
		}
	}

	if len(strErrors) == 0 {
		return nil
	}
	return errors.New(strings.Join(strErrors, "\n"))
}
func (repo *Repository) Mirror(mirrorName string) (*config.Mirror) {
	mirror, ok := repo.mirrorsMaps[mirrorName]
	if ok {
		return &mirror
	}
	return nil
}

func createAuthMethod(mirror config.Mirror) (transport.AuthMethod) {

	var authMethod transport.AuthMethod
	if mirror.AuthMode == "ssh_key" {
		sshKeyPath := "/Users/ralmn/.ssh/id_rsa"

		if mirror.SSHKey != "" {
			sshKeyPath = mirror.SSHKey
		}

		sshKeyAbsPath, err := filepath.Abs(sshKeyPath)
		if err != nil {
			logrus.Error("Failed to find ", sshKeyPath, "... ", err)
			return nil
		}

		publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKeyAbsPath, mirror.Passphrase)
		if err != nil {
			logrus.Error("Failed to create ssh authentification... : ", err)
			return nil
		}

		authMethod = publicKeys

	} else if mirror.AuthMode == "password" {
		password := ssh.Password{
			User:     mirror.User,
			Password: mirror.Password,
		}
		authMethod = &password
	} else {
		logrus.Errorf("Unexpectd auth method %s for mirror %s...", mirror.AuthMode ,mirror)
		return nil
	}

	return authMethod
}

func baresPath() (string) {
	return path.Join("./bares")
}
