// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"errors"
	"os"

	"github.com/spf13/afero"
)

var appFS = afero.NewOsFs()

var (
	ErrCouldNotCreateDirectory = errors.New("belit/fs: could not create directory")
	ErrCouldNotCreateTemporaryFile = errors.New("belit/fs: could not create temporary file")
	ErrFileDoesNotExist = errors.New("belit/fs: file does not exist")
	ErrCouldNotReadFile = errors.New("belit/fs: could not read file")
)

// TODO: add error return
func EnsureDirectory(baseDir string) error {
	err := appFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
		return ErrCouldNotCreateDirectory
	}

	return nil
}

func RemoveFile(name string) error {
	return appFS.Remove(name)
}

func GetTempFile() (string, error) {
	tempDir := afero.GetTempDir(appFS, "belit")

	tempFile, err := afero.TempFile(appFS, tempDir, "belit-")
	if err != nil {
		return "", ErrCouldNotCreateTemporaryFile
	}
	return tempFile.Name(), nil
}

func FileExists(name string) error {
	if _, err := appFS.Stat(name); err != nil {
		return ErrFileDoesNotExist
	}
	return nil
}

func GetFileContents(name string) ([]byte, error) {
	buffer, err := afero.ReadFile(appFS, name)
	if err != nil {
		return []byte{}, ErrCouldNotReadFile
	}
	return buffer, nil
}