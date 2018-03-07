// andrzej lichnerowicz, unlicensed (~public domain)

package sources

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"github.com/unjello/belit/config"
	git "gopkg.in/src-d/go-git.v4"
)

// GitRepo describes components of GitHub remote repository
type GitRepo struct {
	*Source
	protocol string
	site     string
	user     string
	repo     string
	path     string
}

type GitProvider struct{}

// GetName returns name for the provider
func (p *GitProvider) GetName() string { return "git" }

// CanHandle returns true, if so
func (p *GitProvider) CanHandle(repo *Source) (GitRepo, bool) {
	if strings.Contains(repo.uri, "github.com") {
		return describeGitHubRepo(*repo)
	}
	return GitRepo{}, false
}

// GetBasePath returns local path to cloned repository
func (repo *GitRepo) GetBasePath(baseDir string) string {
	return path.Join(baseDir, repo.site, repo.user, repo.repo)
}

// GetIncludePath returns local path deep into repository
// that should be used for includes search
func (repo *GitRepo) GetIncludePath(baseDir string) string {
	return path.Join(baseDir, repo.site, repo.user, repo.repo, repo.path)
}

func (repo *GitRepo) Download(log config.Logger, url string, path string) error {
	// TODO: dodać abstrakcję, zeby móc później mockować
	return fmt.Errorf("Not implemented.")
}

func (repo *GitRepo) getURL() string {
	if repo.protocol == "" {
		return strings.Join([]string{repo.site, repo.user, repo.repo}, "/")
	}

	return fmt.Sprint(repo.protocol, "://", strings.Join([]string{repo.site, repo.user, repo.repo}, "/"))
}

func ensureHTTPSInURL(url string) string {
	i := strings.Index(url, "://")
	if i == -1 {
		return fmt.Sprint("https://", url)
	} else if url[:i] != "https" {
		return fmt.Sprint("https://", url[i+3:])
	}
	return url
}

func describeGitHubRepo(src Source) (GitRepo, bool) {
	var (
		protocol = ""
		site     string
		user     string
		repo     string
		path     = ""
	)

	startIndex := 0
	i := strings.Index(src.uri, "://")
	if i != -1 {
		protocol = src.uri[:i]
		startIndex = i + 3
	}

	bareURL := src.uri[startIndex:]
	bareParts := strings.Split(bareURL, "/")
	if bareParts[0] != "github.com" {
		return GitRepo{}, false
	} else if len(bareParts) < 3 {
		return GitRepo{}, false
	}

	site = bareParts[0]
	user = bareParts[1]
	repo = bareParts[2]

	if len(bareParts) > 3 {
		path = strings.Join(bareParts[3:], "/")
	}

	return GitRepo{&src, protocol, site, user, repo, path}, true
}

// DownloadFromGitHub takes a repo, and if it does not exists
// locally, it clones it to cache directory
func DownloadFromGitHub(baseDir string, url string) error {
	repo, ok := describeGitHubRepo(Source{url, ""})
	if !ok {
		// FIXME: fix
		log.Panic("TODO")
	}

	fullRepoURL := ensureHTTPSInURL(repo.getURL())
	fullBaseDir := repo.GetBasePath(baseDir)

	options := git.CloneOptions{
		URL: fullRepoURL,
	}
	if viper.GetBool("debug") {
		options.Progress = os.Stdout
	}

	if err := ShouldCloneProceed(fullBaseDir); err != nil {
		return nil
	}

	if _, err := git.PlainClone(fullBaseDir, false, &options); err != nil {
		log.Panic(err)
	}

	return nil
}

// ShouldCloneProceed decideds if clone needs to be done
func ShouldCloneProceed(path string) error {
	if _, err := git.PlainOpen(path); err != nil {
		return nil
	}
	return fmt.Errorf("Repository exists")
}
