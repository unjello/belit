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

func DownloadFromGitHub(baseDir string, url string) error {
	repoUrl := ensureHttpsInUrl(url)
	log.WithFields(log.Fields{
		"provider": "github",
		"url":      repoUrl,
		"baseDir":  baseDir,
	}).Info("Starting download process")

	w := log.WithFields(log.Fields{"provider": "git"}).Writer()
	defer w.Close()

	// TODO: Split url to retrieve only repo without subfolders
	// TODO: Make baseDir folder a github.com/<user>/<repo> as in go
	_, err := git.PlainClone(baseDir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	return nil
}
