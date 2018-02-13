// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"os"

	log "github.com/sirupsen/logrus"
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
