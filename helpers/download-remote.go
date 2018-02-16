// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unjello/belit/providers"
)

var AppFS = afero.NewOsFs()

func DownloadRemote(baseDir string, url string) error {
	log.WithFields(log.Fields{
		"folder": baseDir,
		"url":    url,
	}).Info("Preparing to download remote repository")
	// TODO: Add validation of provider
	EnsureDirectory(baseDir)
	srcDir := path.Join(baseDir, "src")
	EnsureDirectory(srcDir)
	return providers.DownloadFromGitHub(srcDir, url)
}
