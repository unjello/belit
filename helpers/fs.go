// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func EnsureDirectory(baseDir string) {
	log.WithFields(log.Fields{
		"folder": baseDir,
	}).Debug("Ensure folder exists")
	err := AppFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
		log.WithFields(log.Fields{
			"folder": baseDir,
			"error":  err,
		}).Errorf("Error creating folder")
	}
}

func RemoveFile(name string) error {
	log.WithFields(log.Fields{
		"file": name,
	}).Debug("Remove file")
	return AppFS.Remove(name)
}

func GetTempFile() (string, error) {
	tempDir := afero.GetTempDir(AppFS, "belit")
	l := log.WithFields(log.Fields{
		"folder": tempDir,
	})
	l.Debug("Created temporary folder")

	tempFile, err := afero.TempFile(AppFS, tempDir, "belit-")
	if err != nil {
		l.Error("Error creating temporary file")
		return "", err
	}
	log.WithFields(log.Fields{
		"file": tempFile.Name(),
	}).Debug("Created temporary file")

	return tempFile.Name(), nil
}

func FileExists(name string) error {
	l := log.WithFields(log.Fields{
		"file": name,
	})
	l.Info("Check if file exists")
	if _, err := AppFS.Stat("/path/to/whatever"); err != nil {
		l.Info("File does not exist")
		return fmt.Errorf("File does not exist")
	}
	return nil
}
