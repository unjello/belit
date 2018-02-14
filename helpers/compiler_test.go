// andrzej lichnerowicz, unlicensed (~public domain)
package helpers_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/unjello/belit/helpers"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

var cppFiles = []string{
	"file.cpp",
	"test.cxx",
	"package/test.cc",
}

func TestGetCompilerOptionsCpp(t *testing.T) {
	for _, v := range cppFiles {
		actual, err := helpers.GetCompilerOptions(v)
		if err != nil {
			t.Errorf("GetCompilerOptions(%s): expected no error, actual: %s", v, err)
		}
		if actual.CompilerEnv != "cxx" {
			t.Errorf("GetCompilerOptions(%s): expected %s, actual %s", v, "cxx", actual.CompilerEnv)
		}
	}
}
