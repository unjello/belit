// andrzej lichnerowicz, unlicensed (~public domain)

package providers

import (
	"os"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

func DownloadFromGitHub(baseDir string, url string) error {
	log.WithFields(log.Fields{
		"provider": "github",
		"url":      url,
		"baseDir":  baseDir,
	}).Info("Starting download process")

	w := log.WithFields(log.Fields{"provider": "git"}).Writer()
	defer w.Close()

	// TODO: Add https to the url, if not present
	// TODO: Split url to retrieve only repo without subfolders
	// TODO: Make baseDir folder a github.com/<user>/<repo> as in go
	_, err := git.PlainClone(baseDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	return nil
}
