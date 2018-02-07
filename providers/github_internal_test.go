// andrzej lichnerowicz, unlicensed (~public domain)
package providers

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

var httpsInURLData = []struct {
	url      string
	expected string
}{
	{"github.com/", "https://github.com/"},
	{"github.com/path/to/repo", "https://github.com/path/to/repo"},
	{"http://github.com/user", "https://github.com/user"},
	{"https://github.com/", "https://github.com/"},
}

func TestEnsureHttpsInUrl(t *testing.T) {
	for _, v := range httpsInURLData {
		actual := ensureHttpsInUrl(v.url)
		if actual != v.expected {
			t.Errorf("ensureHttpsInUrl(%s): expected %s, actual %s", v.url, v.expected, actual)
		}
	}
}
