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
	src      Source
	expected GitRepo
}{
	{Source{"http://github.com/user/repo/", "" }, GitRepo{&Source{"http://github.com/user/repo/", "" }, "http", "github.com", "user", "repo", ""}},
	{Source{"github.com/user/repo/","" }, GitRepo{&Source{"github.com/user/repo/","" },"", "github.com", "user", "repo", ""}},
	{Source{"github.com/user/repo/folder/", "" },GitRepo{&Source{"github.com/user/repo/folder/", "" },"", "github.com", "user", "repo", "folder/"}},
	{Source{"https://github.com/user/repo/folder/deeper/", "" },GitRepo{&Source{"https://github.com/user/repo/folder/deeper/", "" }, "https", "github.com", "user", "repo", "folder/deeper/"}},
}

func TestGetGitRepo(t *testing.T) {
	for _, v := range gitRepoData {
		actual, ok := describeGitHubRepo(v.src)

		assert.True(t, ok)
		assert.Equal(t, v.expected, actual)
	}
}

var gitRepoInvalidData = []Source{
	Source{"http://github.com/user", "" },
	Source{"github.com/user", "" },
	Source{"github.com/", "" },
	Source{"gitlab.com", "" },
	Source{"gitlab.com/user/repo/", "" },
	Source{"http://gitlab.com/user/", "" },
	Source{"other.com", "" },
}

func TestGetGitRepoInvalid(t *testing.T) {
	for _, src := range gitRepoInvalidData {
		_, ok := describeGitHubRepo(src)

		assert.False(t, ok)
	}
}

var gitRepoURLData = []struct {
	repo     GitRepo
	expected string
}{
	{GitRepo{&Source{ "", ""}, "http", "github.com", "user", "repo", ""}, "http://github.com/user/repo"},
	{GitRepo{&Source{"", ""}, "", "github.com", "user", "repo", ""}, "github.com/user/repo"},
	{GitRepo{&Source{"", ""}, "", "github.com", "user", "repo", "folder/"}, "github.com/user/repo"},
	{GitRepo{&Source{"", ""}, "https", "github.com", "user", "repo", "folder/deeper/"}, "https://github.com/user/repo"},
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
		repo, ok := describeGitHubRepo(Source{v.url, ""} )
		assert.True(t, ok)

		actual := repo.GetIncludePath(baseDir)
		assert.Equal(t, v.expected, actual)
	}
}
