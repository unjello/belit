// andrzej lichnerowicz, unlicensed (~public domain)
package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gitRepoData = []struct {
	src      string
	protocol string
	site string
	user string
	repo string
	path string
}{
	{ "http://github.com/user/repo/", "http", "github.com", "user", "repo", "" },
	{ "github.com/user/repo/", "", "github.com", "user", "repo", "" },
	{ "github.com/user/repo/folder/", "", "github.com", "user", "repo", "folder/" },
	{ "https://github.com/user/repo/folder/deeper/", "https", "github.com", "user", "repo", "folder/deeper/" },
}

func TestGetGitRepo(t *testing.T) {
	for _, v := range gitRepoData {
		protocolA, siteA, userA, repoA, pathA , ok := describeGitHubRepo(v.src)

		assert.True(t, ok)
		assert.Equal(t, v.protocol, protocolA)
		assert.Equal(t, v.site, siteA)
		assert.Equal(t, v.user, userA)
		assert.Equal(t, v.repo, repoA)
		assert.Equal(t, v.path, pathA)
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
	for _, src := range gitRepoInvalidData {
		_, _, _, _, _, ok := describeGitHubRepo(src)

		assert.False(t, ok)
	}
}

var gitRepoURLData = []struct {
	repo     GitProvider
	expected string
}{
	{GitProvider{"http", "github.com", "user", "repo", ""}, "http://github.com/user/repo"},
	{GitProvider{"", "github.com", "user", "repo", ""}, "github.com/user/repo"},
	{GitProvider{"", "github.com", "user", "repo", "folder/"}, "github.com/user/repo"},
	{GitProvider{"https", "github.com", "user", "repo", "folder/deeper/"}, "https://github.com/user/repo"},
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
		protocol, site, user, repo, path, ok := describeGitHubRepo(v.url)
		assert.True(t, ok)

		p := &GitProvider{protocol, site, user, repo, path}
		actual := p.GetIncludePath(baseDir)
		assert.Equal(t, v.expected, actual)
	}
}
