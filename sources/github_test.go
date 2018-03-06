// andrzej lichnerowicz, unlicensed (~public domain)
package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		actual := ensureHTTPSInURL(v.url)

		assert.Equal(t, v.expected, actual)
	}
}

var gitRepoData = []struct {
	url      string
	expected GitRepo
}{
	{"http://github.com/user/repo/", GitRepo{"http", "github.com", "user", "repo", ""}},
	{"github.com/user/repo/", GitRepo{"", "github.com", "user", "repo", ""}},
	{"github.com/user/repo/folder/", GitRepo{"", "github.com", "user", "repo", "folder/"}},
	{"https://github.com/user/repo/folder/deeper/", GitRepo{"https", "github.com", "user", "repo", "folder/deeper/"}},
}

func TestGetGitRepo(t *testing.T) {
	for _, v := range gitRepoData {
		actual, err := GetGitRepo(v.url)

		assert.Nil(t, err)
		assert.Equal(t, v.expected, actual)
	}
}

var gitRepoInvalidData = []string{
	"http://github.com/user",
	"github.com/user",
	"github.com/",
	"gitlab.com",
	"gitlab.com/user/repo/",
	"http://gitlab.com/user/",
	"other.com",
}

func TestGetGitRepoInvalid(t *testing.T) {
	for _, url := range gitRepoInvalidData {
		_, err := GetGitRepo(url)

		assert.NotNil(t, err)
	}
}

var gitRepoURLData = []struct {
	repo     GitRepo
	expected string
}{
	{GitRepo{"http", "github.com", "user", "repo", ""}, "http://github.com/user/repo"},
	{GitRepo{"", "github.com", "user", "repo", ""}, "github.com/user/repo"},
	{GitRepo{"", "github.com", "user", "repo", "folder/"}, "github.com/user/repo"},
	{GitRepo{"https", "github.com", "user", "repo", "folder/deeper/"}, "https://github.com/user/repo"},
}

func TestGitRepoGetUrl(t *testing.T) {
	for _, v := range gitRepoURLData {
		url := v.repo.getURL()

		assert.Equal(t, v.expected, url)
	}
}

func TestGitRepoGetBasePath(t *testing.T) {
	const baseDir = "/home/xxx/.belit/src"
	const expected = "/home/xxx/.belit/src/github.com/user/repo"

	for _, v := range gitRepoURLData {
		path := v.repo.GetBasePath(baseDir)

		assert.Equal(t, expected, path)
	}
}

func TestGitRepoGetIncludePath(t *testing.T) {
	const baseDir = "/home/xxx/.belit/src"
	var gitRepoData = []struct {
		url      string
		expected string
	}{
		{"https://github.com/catchorg/Catch2/single_include/", "/home/xxx/.belit/src/github.com/catchorg/Catch2/single_include"},
		{"https://github.com/catchorg/Catch2/single_include/even/more", "/home/xxx/.belit/src/github.com/catchorg/Catch2/single_include/even/more"},
	}

	for _, v := range gitRepoData {
		repo, err := GetGitRepo(v.url)
		assert.Nil(t, err)

		actual := repo.GetIncludePath(baseDir)
		assert.Equal(t, v.expected, actual)
	}
}
