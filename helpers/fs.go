// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unjello/belit/config"
)

func EnsureDirectory(baseDir string) {
	log := config.GetConfig().Log
	log.WithFields(logrus.Fields{
		"folder": baseDir,
	}).Debug("Ensure folder exists")
	err := AppFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
		log.WithFields(logrus.Fields{
			"folder": baseDir,
			"error":  err,
		}).Errorf("Error creating folder")
	}
}

func RemoveFile(name string) error {
	log := config.GetConfig().Log
	log.WithFields(logrus.Fields{
		"file": name,
	}).Debug("Remove file")
	return AppFS.Remove(name)
}

func GetTempFile() (string, error) {
	tempDir := afero.GetTempDir(AppFS, "belit")
	log := config.GetConfig().Log
	l := log.WithFields(logrus.Fields{
		"folder": tempDir,
	})
	l.Debug("Created temporary folder")

	tempFile, err := afero.TempFile(AppFS, tempDir, "belit-")
	if err != nil {
		l.Error("Error creating temporary file")
		return "", err
	}
	log.WithFields(logrus.Fields{
		"file": tempFile.Name(),
	}).Debug("Created temporary file")

	return tempFile.Name(), nil
}

func FileExists(name string) error {
	log := config.GetConfig().Log
	l := log.WithFields(logrus.Fields{
		"file": name,
	})
	l.Debug("Check if file exists")
	if _, err := AppFS.Stat(name); err != nil {
		l.Info("File does not exist")
		return fmt.Errorf("File does not exist")
	}
	l.Info("File exists")
	return nil
}
