// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

func TestGetCompilerOptions(t *testing.T) {
	var compilerFiles = []struct {
		file     string
		expected string
	}{
		{"file.cpp", "cxx"},
		{"test.cxx", "cxx"},
		{"package/test.cc", "cxx"},
		{"/outer/inner/file.c++", "cxx"},
		{"/test.pcc", "cxx"},
		{"file.c", "cc"},
		{"/folder/file.c", "cc"},
	}

	for _, v := range compilerFiles {
		actual, err := GetCompilerOptions(v.file)

		assert.Nil(t, err)
		assert.Equal(t, actual.CompilerEnv, v.expected)
	}
}

func TestGetUnknownCompilerOptions(t *testing.T) {
	var unknownCompilerFiles = []string{
		"file.txt",
		"file.go",
		"/folder/file.csv",
	}

	for _, v := range unknownCompilerFiles {
		_, err := GetCompilerOptions(v)
		assert.NotNil(t, err)
	}
}
