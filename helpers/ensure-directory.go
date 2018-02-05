// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func EnsureDirectory(baseDir string) {
	log.WithFields(log.Fields{
		"baseDir": baseDir,
	}).Info("Ensuring folder exists")
	err := AppFS.MkdirAll(baseDir, os.ModeDir|0775)
	if err != nil {
		log.WithFields(log.Fields{
			"baseDir": baseDir,
			"error":   err,
		}).Errorf("Error creating folder")
	}
}
