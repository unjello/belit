package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestEnsureHttps(t *testing.T) {
	var httpsInURLData = []struct {
		url      string
		expected string
	}{
		{"github.com/", "https://github.com/"},
		{"github.com/path/to/repo", "https://github.com/path/to/repo"},
		{"http://github.com/user", "https://github.com/user"},
		{"https://github.com/", "https://github.com/"},
	}

	for _, v := range httpsInURLData {
		actual := ensureHTTPS(v.url)

		assert.Equal(t, v.expected, actual)
	}
}
