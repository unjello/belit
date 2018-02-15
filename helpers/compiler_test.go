// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

var cppFiles = []string{
	"file.cpp",
	"test.cxx",
	"package/test.cc",
	"/outer/inner/file.c++",
	"/test.pcc",
}

func TestGetCompilerOptionsCpp(t *testing.T) {
	for _, v := range cppFiles {
		actual, err := GetCompilerOptions(v)
		if err != nil {
			t.Errorf("GetCompilerOptions(%s): expected no error, actual: %s", v, err)
		}
		if actual.CompilerEnv != "cxx" {
			t.Errorf("GetCompilerOptions(%s): expected %s, actual %s", v, "cxx", actual.CompilerEnv)
		}
	}
}
