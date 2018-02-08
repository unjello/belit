// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

func ensureHttpsInUrl(url string) string {
	i := strings.Index(url, "://")
	l := log.WithFields(log.Fields{
		"url": url,
	})
	if i == -1 {
		l.Info("No prefix detected. Adding https://.")
		return fmt.Sprint("https://", url)
	} else if url[:i] != "https" {
		l.Info("Wrong prefix detected. Changing to https://.")
		return fmt.Sprint("https://", url[i+3:])
	}
	return url
}

func getGitHubRepoUrl(url string) (string, error) {
	i := strings.Index(url, "://")
	l := log.WithFields(log.Fields{
		"url": url,
	})

	startIndex := 0
	if i != -1 {
		startIndex = i + 3
	}

	bareUrl := url[startIndex:]
	bareParts := strings.Split(bareUrl, "/")
	if bareParts[0] != "github.com" {
		l.Warn("URL doesn't seem to point to GitHub")
		return url, fmt.Errorf("URL doesn't seem to point to GitHub")
	} else if len(bareParts) < 3 {
		l.Warn("URL too short. GitHub repo follows format: github.com/<user>/<repo>/")
		return url, fmt.Errorf("URL too short. GitHub repo follows format: github.com/<user>/<repo>/")
	}

	newUrl := fmt.Sprint(url[:startIndex], strings.Join(bareParts[:3], "/"))
	log.WithFields(log.Fields{
		"oldUrl": url,
		"newUrl": newUrl,
	}).Info("Extracted GitHub repository URL")
	return newUrl, nil
}

func DownloadFromGitHub(baseDir string, url string) error {
	repoUrl, errr := getGitHubRepoUrl(url)
	if errr != nil {
		panic(errr)
	}
	fullRepoUrl := ensureHttpsInUrl(repoUrl)
	log.WithFields(log.Fields{
		"provider": "github",
		"url":      fullRepoUrl,
		"baseDir":  baseDir,
	}).Info("Starting download process")

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
