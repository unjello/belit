// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"path"

	"github.com/spf13/afero"
	"github.com/unjello/belit/providers"
)

var AppFS = afero.NewOsFs()

func DownloadRemote(baseDir string, url string) error {
	// TODO: Add validation of provider
	EnsureDirectory(baseDir)
	srcDir := path.Join(baseDir, "src")
	EnsureDirectory(srcDir)
	return providers.DownloadFromGitHub(srcDir, url)
}
