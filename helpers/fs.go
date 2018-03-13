// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

var appFS = afero.NewOsFs()

// TODO: add error return
func EnsureDirectory(baseDir string) {
	err := appFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
	}
}

func RemoveFile(name string) error {
	return appFS.Remove(name)
}

func GetTempFile() (string, error) {
	tempDir := afero.GetTempDir(appFS, "belit")

	tempFile, err := afero.TempFile(appFS, tempDir, "belit-")
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func FileExists(name string) error {
	if _, err := appFS.Stat(name); err != nil {
		return fmt.Errorf("File does not exist")
	}
	return nil
}

func GetFileContents(name string) ([]byte, error) {
	return afero.ReadFile(appFS, name)
}