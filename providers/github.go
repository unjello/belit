// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"fmt"
	"os"
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
		l.Warn(warnNotGithubUrl)
		return url, fmt.Errorf(warnNotGithubUrl)
	} else if len(bareParts) < 3 {
		l.Warn(warnUrlTooShort)
		return url, fmt.Errorf(warnUrlTooShort)
	}

	newUrl := fmt.Sprint(url[:startIndex], strings.Join(bareParts[:3], "/"))
	log.WithFields(log.Fields{
		"oldUrl": url,
		"newUrl": newUrl,
	}).Info(infoExtractedUrl)
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
