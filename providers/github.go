// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unjello/belit/config"
	git "gopkg.in/src-d/go-git.v4"
)

// GitRepo describes components of GitHub remote repository
type GitRepo struct {
	protocol string
	site     string
	user     string
	repo     string
	path     string
}

// GetBasePath returns local path to cloned repository
func (repo *GitRepo) GetBasePath(baseDir string) string {
	return path.Join(baseDir, repo.site, repo.user, repo.repo)
}

// GetIncludePath returns local path deep into repository
// that should be used for includes search
func (repo *GitRepo) GetIncludePath(baseDir string) string {
	return path.Join(baseDir, repo.site, repo.user, repo.repo, repo.path)
}

// GetGitRepo breaks down URI into separate components
func GetGitRepo(log config.Logger, url string) (GitRepo, error) {
	var (
		protocol = ""
		site     string
		user     string
		repo     string
		path     = ""
	)

	_log := log.WithFields(logrus.Fields{
		"url": url,
	})

	startIndex := 0
	i := strings.Index(url, "://")
	if i != -1 {
		protocol = url[:i]
		startIndex = i + 3
	}

	bareURL := url[startIndex:]
	bareParts := strings.Split(bareURL, "/")
	if bareParts[0] != "github.com" {
		err := fmt.Errorf("URL doesn't seem to point to GitHub")
		_log.Warn(err)
		return GitRepo{}, err
	} else if len(bareParts) < 3 {
		err := fmt.Errorf("URL too short. GitHub repo follows format: github.com/<user>/<repo>/")
		_log.Warn(err)
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

func (repo *GitRepo) getURL(log config.Logger) string {
	var newURL string
	if repo.protocol == "" {
		newURL = strings.Join([]string{repo.site, repo.user, repo.repo}, "/")
	} else {
		newURL = fmt.Sprint(repo.protocol, "://", strings.Join([]string{repo.site, repo.user, repo.repo}, "/"))
	}

	log.WithFields(logrus.Fields{
		"protocol": repo.protocol,
		"site":     repo.site,
		"user":     repo.user,
		"repo":     repo.repo,
		"newUrl":   newURL,
	}).Info("Extracted GitHub repository URL")

	return newURL
}

func ensureHTTPSInURL(log config.Logger, url string) string {
	i := strings.Index(url, "://")
	l := log.WithFields(logrus.Fields{
		"url": url,
	})

	if i == -1 {
		newURL := fmt.Sprint("https://", url)
		l.WithFields(logrus.Fields{
			"new": newURL,
		}).Debug("No prefix detected. Adding https://.")

		return newURL
	} else if url[:i] != "https" {
		newURL := fmt.Sprint("https://", url[i+3:])
		l.WithFields(logrus.Fields{
			"new": newURL,
		}).Debug("Wrong prefix detected. Changing to https://.")

		return newURL
	}

	return url
}

// DownloadFromGithub takes a repo, and if it does not exists
// locally, it clones it to cache directory
func DownloadFromGitHub(log config.Logger, baseDir string, url string) error {
	repo, err := GetGitRepo(log, url)
	if err != nil {
		log.Panic(err)
	}

	fullRepoURL := ensureHTTPSInURL(log, repo.getURL(log))
	fullBaseDir := repo.GetBasePath(baseDir)

	log.WithFields(logrus.Fields{
		"provider": "github",
		"url":      fullRepoURL,
		"baseDir":  fullBaseDir,
	}).Info("Starting download process")

	options := git.CloneOptions{
		URL: fullRepoURL,
	}
	if viper.GetBool("debug") {
		options.Progress = os.Stdout
	}

	if err := CheckFolderBeforeGitClone(log, fullBaseDir); err != nil {
		log.WithFields(logrus.Fields{
			"path": fullBaseDir,
		}).Warn("Folder already contain valid Git repository")
		return nil
	}

	if _, err := git.PlainClone(fullBaseDir, false, &options); err != nil {
		log.Panic(err)
	}

	return nil
}

// CheckFolderBeforeGitClone decideds if clone needs to be done
func CheckFolderBeforeGitClone(log config.Logger, path string) error {
	_log := log.WithFields(logrus.Fields{
		"path": path,
	})
	_log.Debug("Check if path contains valid git repository")

	if _, err := git.PlainOpen(path); err != nil {
		_log.Debug("Folder does not contain git repository")
		return nil
	}

	_log.Debug("Folder does contain valid git repository")
	return fmt.Errorf("Repository exists")
}
