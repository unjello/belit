// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/unjello/belit/config"
)

type CompilerOptions struct {
	CompilerEnv        string
	CompilerOptionsEnv string
}

func GetCompilerOptions(filename string) (CompilerOptions, error) {
	ext := filepath.Ext(filename)
	log := config.GetConfig().Log
	log.WithFields(logrus.Fields{
		"file": filename,
		"ext":  ext,
	}).Info("Looking for matching compiler")
	switch ext {
	case ".c":
		return CompilerOptions{"cc", "ccopts"}, nil
	case ".cc":
		fallthrough
	case ".c++":
		fallthrough
	case ".cpp":
		fallthrough
	case ".pcc":
		fallthrough
	case ".cxx":
		return CompilerOptions{"cxx", "cxxopts"}, nil
	default:
		return CompilerOptions{}, fmt.Errorf("Unknown language. Cannot match compiler")
	}
}
