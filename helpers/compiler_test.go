// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

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

func TestGetCompilerOptions(t *testing.T) {
	for _, v := range compilerFiles {
		actual, err := GetCompilerOptions(v.file)
		if err != nil {
			t.Errorf("GetCompilerOptions(%s): expected no error, actual: %s", v, err)
		}
		if actual.CompilerEnv != v.expected {
			t.Errorf("GetCompilerOptions(%s): expected %s, actual %s", v.file, v.expected, actual.CompilerEnv)
		}
	}
}
