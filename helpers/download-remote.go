// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unjello/belit/config"
	"github.com/unjello/belit/sources"
)

var AppFS = afero.NewOsFs()

func DownloadRemote(baseDir string, url string) error {
	log := config.GetConfig().Log
	log.WithFields(logrus.Fields{
		"folder": baseDir,
		"url":    url,
	}).Info("Preparing to download remote repository")
	// TODO: Add validation of provider
	EnsureDirectory(baseDir)
	srcDir := path.Join(baseDir, "src")
	EnsureDirectory(srcDir)

	provider := sources.GitProvider{}
	if provider.CanHandle(url) {
		return provider.Download(srcDir, false)
	}

	return fmt.Errorf("Unknown remote type")
}
