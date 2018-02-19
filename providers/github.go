// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
)

func ensureHttpsInUrl(url string) string {
	i := strings.Index(url, "://")
	l := log.WithFields(log.Fields{
		"url": url,
	})
	if i == -1 {
		newUrl := fmt.Sprint("https://", url)
		l.WithFields(log.Fields{
			"new": newUrl,
		}).Debug("No prefix detected. Adding https://.")
		return newUrl
	} else if url[:i] != "https" {
		newUrl := fmt.Sprint("https://", url[i+3:])
		l.WithFields(log.Fields{
			"new": newUrl,
		}).Debug("Wrong prefix detected. Changing to https://.")
		return newUrl
	}
	return url
}

func (repo *GitRepo) getUrl() string {
	l := log.WithFields(log.Fields{
		"protocol": repo.protocol,
		"site":     repo.site,
		"user":     repo.user,
		"repo":     repo.repo,
	})

	var newUrl string
	if repo.protocol == "" {
		newUrl = strings.Join([]string{repo.site, repo.user, repo.repo}, "/")
	} else {
		newUrl = fmt.Sprint(repo.protocol, "://", strings.Join([]string{repo.site, repo.user, repo.repo}, "/"))
	}
	l.WithFields(log.Fields{
		"newUrl": newUrl,
	}).Info("Extracted GitHub repository URL")
	return newUrl
}

func (repo *GitRepo) getBasePath(baseDir string) string {
	return path.Join(baseDir, repo.site, repo.user, repo.repo)
}

type GitRepo struct {
	protocol string
	site     string
	user     string
	repo     string
	path     string
}

func getGitRepo(url string) (GitRepo, error) {
	var (
		protocol = ""
		site     string
		user     string
		repo     string
		path     = ""
	)

	l := log.WithFields(log.Fields{
		"url": url,
	})

	startIndex := 0
	i := strings.Index(url, "://")
	if i != -1 {
		protocol = url[:i]
		startIndex = i + 3
	}

	bareUrl := url[startIndex:]
	bareParts := strings.Split(bareUrl, "/")
	if bareParts[0] != "github.com" {
		err := fmt.Errorf("URL doesn't seem to point to GitHub")
		l.Warn(err)
		return GitRepo{}, err
	} else if len(bareParts) < 3 {
		err := fmt.Errorf("URL too short. GitHub repo follows format: github.com/<user>/<repo>/")
		l.Warn(err)
		return GitRepo{}, err
	}

	site = bareParts[0]
	user = bareParts[1]
	repo = bareParts[2]

	if len(bareParts) > 3 {
		path = strings.Join(bareParts[3:], "/")
	}

	return GitRepo{protocol, site, user, repo, path}, nil
}

func DownloadFromGitHub(baseDir string, url string) error {
	repo, errr := getGitRepo(url)
	if errr != nil {
		panic(errr)
	}
	fullRepoUrl := ensureHttpsInUrl(repo.getUrl())
	fullBaseDir := repo.getBasePath(baseDir)
	log.WithFields(log.Fields{
		"provider": "github",
		"url":      fullRepoUrl,
		"baseDir":  fullBaseDir,
	}).Info("Starting download process")

	w := log.WithFields(log.Fields{"provider": "git"}).Writer()
	defer w.Close()

	options := git.CloneOptions{
		URL: fullRepoUrl,
	}
	if viper.GetBool("debug") {
		options.Progress = os.Stdout
	}

	if err := CheckFolderBeforeGitClone(fullBaseDir); err != nil {
		log.WithFields(log.Fields{
			"path": fullBaseDir,
		}).Warn("Folder already contain valid Git repository")
		return nil
	}
	_, err := git.PlainClone(fullBaseDir, false, &options)

	if err != nil {
		panic(err)
	}

	return nil
}

func CheckFolderBeforeGitClone(path string) error {
	l := log.WithFields(log.Fields{
		"path": path,
	})
	l.Debug("Check if path contains valid git repository")
	if _, err := git.PlainOpen(path); err != nil {
		l.Debug("Folder does not contain git repository")
		return nil
	}

	l.Debug("Folder does contain valid git repository")
	return fmt.Errorf("Repository exists")
}
