package helpers

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestGetSources(t *testing.T) {
	var sources = []struct {
		text string
		expected []SourceInfo
	}{
		{`#include /* github.com/catchorg/Catch2/single_include/ */ "catch.hpp"`,
		[]SourceInfo{
					SourceInfo{ "github.com/catchorg/Catch2/single_include/", "catch.hpp" },
					}},
		{`#include /* github.com/catchorg/Catch2/single_include/ */ "catch.hpp"
#include/*https://github.com/nothings/stb*/"stb.h"`,
		[]SourceInfo{
					SourceInfo{ "github.com/catchorg/Catch2/single_include/", "catch.hpp" },
					SourceInfo{ "https://github.com/nothings/stb", "stb.h" },
					}},
		{`#include "catch.hpp"
#include <cassert>
#include/*https://github.com/nothings/stb*/"stb.h"`,
			[]SourceInfo{
				SourceInfo{ "https://github.com/nothings/stb", "stb.h" },
			}},
	}

	for _, s := range sources {
		actual, err := GetSourcesFromBuffer([]byte(s.text))
		assert.Nil(t, err)
		for i, a := range actual {
			assert.Equal(t, s.expected[i].HeaderName, a.HeaderName)
			assert.Equal(t, s.expected[i].RepositoryPath, a.RepositoryPath)
		}
	}
}