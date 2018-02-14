// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type CompilerOptions struct {
	CompilerEnv        string
	CompilerOptionsEnv string
}

func GetCompilerOptions(filename string) (CompilerOptions, error) {
	ext := filepath.Ext(filename)
	log.WithFields(log.Fields{
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
