// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

const (
	infoExtractedUrl         = "Extracted GitHub repository URL"
	infoStartingDownload     = "Starting download process"
	infoNoPrefixDetected     = "No prefix detected. Adding https://."
	infoWrongPrefixdDetected = "Wrong prefix detected. Changing to https://."
	warnNotGithubUrl         = "URL doesn't seem to point to GitHub"
	warnUrlTooShort          = "URL too short. GitHub repo follows format: github.com/<user>/<repo>/"
)

func ensureHttpsInUrl(url string) string {
	i := strings.Index(url, "://")
	l := log.WithFields(log.Fields{
		"url": url,
	})
	if i == -1 {
		l.Info(infoNoPrefixDetected)
		return fmt.Sprint("https://", url)
	} else if url[:i] != "https" {
		l.Info(infoWrongPrefixdDetected)
		return fmt.Sprint("https://", url[i+3:])
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
	}).Info(infoExtractedUrl)
	return newUrl
}

func (repo *GitRepo) getBasePath(baseDir string) string {
	return path.Join(baseDir, "src", repo.site, repo.user, repo.repo)
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
		l.Warn(warnNotGithubUrl)
		return GitRepo{}, fmt.Errorf(warnNotGithubUrl)
	} else if len(bareParts) < 3 {
		l.Warn(warnUrlTooShort)
		return GitRepo{}, fmt.Errorf(warnUrlTooShort)
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
	log.WithFields(log.Fields{
		"provider": "github",
		"url":      fullRepoUrl,
		"baseDir":  baseDir,
	}).Info(infoStartingDownload)

	w := log.WithFields(log.Fields{"provider": "git"}).Writer()
	defer w.Close()

	// TODO: Make baseDir folder a github.com/<user>/<repo> as in go
	_, err := git.PlainClone(baseDir, false, &git.CloneOptions{
		URL:      fullRepoUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	return nil
}
