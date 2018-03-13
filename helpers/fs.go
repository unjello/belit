// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

func EnsureDirectory(baseDir string) {
	err := AppFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
	}
}

func RemoveFile(name string) error {
	return AppFS.Remove(name)
}

func GetTempFile() (string, error) {
	tempDir := afero.GetTempDir(AppFS, "belit")

	tempFile, err := afero.TempFile(AppFS, tempDir, "belit-")
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func FileExists(name string) error {
	if _, err := AppFS.Stat(name); err != nil {
		return fmt.Errorf("File does not exist")
	}
	return nil
}

func GetFileContents(name string) ([]byte, error) {
	return afero.ReadFile(AppFS, name)
}