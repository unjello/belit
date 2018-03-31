// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"path/filepath"
)

// CompilerOptions holds names of environments to look at
// based on file type (extension)
type CompilerOptions struct {
	CompilerEnv        string
	CompilerOptionsEnv string
}

// GetCompilerOptions returns set of names to corresponding options
// based on file type
func GetCompilerOptions(filename string) (CompilerOptions, error) {
	ext := filepath.Ext(filename)
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
