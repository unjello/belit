package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractModelineOptions(t *testing.T) {
	var sources = []struct {
		source   string
		expected map[string]string
	}{
		{`// belit: cxx=g++-7 cxxopts="-O2 -std=c++17"`, map[string]string{"cxx": "g++-7", "cxxopts": "-O2 -std=c++17"}},
		{`// belit: cxx="g++-7" cxxopts="-O2 -std=c++17"`, map[string]string{"cxx": "g++", "cxxopts": "-O2 -std=c++17"}},
		{`/* belit: cc=gcc ccopts=-O0 */`, map[string]string{"cc": "gcc", "ccopts": "-O0"}},
		{`/* belit: cxx="clang-8.0" cxxopts="-O0 -D_WIN32" */`, map[string]string{"cxx": "clang-8.0", "cxxopts": "-O0 -D_WIN32"}},
	}

	for _, e := range sources {
		options, err := ExtractModelineOptions(e.source)
		assert.Nil(t, err)
		for k, v := range e.expected {
			assert.Equal(t, v, options[k])
		}
	}
}

func TestExtractModeline(t *testing.T) {
	var sources = []struct {
		source   string
		expected string
	}{
		{`// belit: cxx=g++-7 cxxopts="-O2 -std=c++17"`, `cxx=g++-7 cxxopts="-O2 -std=c++17"`},
		{`// belit: cxx="g++-7" cxxopts="-O2 -std=c++17"`, `cxx="g++-7" cxxopts="-O2 -std=c++17"`},
		{`/* belit: cc=gcc ccopts=-O0 */`, `cc=gcc ccopts=-O0`},
		{`/* belit: cxx="clang-8.0" cxxopts="-O0 -D_WIN32" */`, `cxx="clang-8.0" cxxopts="-O0 -D_WIN32"`},
	}

	for _, e := range sources {
		modeline, err := extractModeline(e.source)
		assert.Nil(t, err)
		assert.Equal(t, modeline, e.expected)
	}
}
