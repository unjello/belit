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

var githubRepoUrlData = []struct {
	url      string
	expected GitRepo
}{
	{"http://github.com/user/repo/", GitRepo{"http", "github.com", "user", "repo", ""}},
	{"github.com/user/repo/", GitRepo{"", "github.com", "user", "repo", ""}},
	{"github.com/user/repo/folder/", GitRepo{"", "github.com", "user", "repo", "folder/"}},
	{"https://github.com/user/repo/folder/deeper/", GitRepo{"https", "github.com", "user", "repo", "folder/deeper/"}},
}

func TestGetGitRepo(t *testing.T) {
	for _, v := range githubRepoUrlData {
		actual, err := getGitRepo(v.url)
		if err != nil {
			t.Errorf("getGitRepo(%s): expected success, actual error: %s", v.url, err)
		}
		if actual != v.expected {
			t.Errorf("getGitRepo(%s): expected %+v, actual %+v", v.url, v.expected, actual)
		}
	}
}

var githubRepoUrlInvalidData = []string{
	"http://github.com/user",
	"github.com/user",
	"github.com/",
	"gitlab.com",
	"gitlab.com/user/repo/",
	"http://gitlab.com/user/",
	"other.com",
}

func TestGetGitHubRepoUrlInvalid(t *testing.T) {
	for _, url := range githubRepoUrlInvalidData {
		_, err := getGitHubRepoUrl(url)
		if err == nil {
			t.Errorf("getGitHubRepoUrl(%s): expected error, actual success", url)
		}
	}
}
