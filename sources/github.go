// andrzej lichnerowicz, unlicensed (~public domain)

package sources

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

// GitRepo describes components of GitHub remote repository
type GitProvider struct {
	protocol string
	site     string
	user     string
	repo     string
	folder   string
}

// GetName returns name for the provider
func (p *GitProvider) GetName() string { return "git" }

// CanHandle returns true, if so
func (p *GitProvider) CanHandle(uri string) bool {
	if strings.Contains(uri, "github.com") {
		protocol, site, user, repo, folder, ok := describeGitHubRepo(uri)
		if ok {
			p.protocol = protocol
			p.site = site
			p.user = user
			p.repo = repo
			p.folder = folder
			return true
		}
	}
	return false
}

// Download fetches remote repository
func (p *GitProvider) Download(path string, debug bool, handlers ...interface{}) error {
	uri := ensureHTTPS(p.getURL())
	fullPath := p.GetBasePath(path)

	var handler remoteHandler = gitV4Handler{}
	if len(handlers) == 1 {
		handler = handlers[0].(remoteHandler)
	}

	return DownloadFromGitHub(handler, uri, fullPath, debug)
}

func init() {
	RegisterProvider(&GitProvider{})
}

// GetBasePath returns local path to cloned repository
func (p *GitProvider) GetBasePath(baseDir string) string {
	return path.Join(baseDir, p.site, p.user, p.repo)
}

// GetIncludePath returns local path deep into repository
// that should be used for includes search
func (p *GitProvider) GetIncludePath(baseDir string) string {
	return path.Join(baseDir, p.site, p.user, p.repo, p.folder)
}

func (p *GitProvider) getURL() string {
	if p.protocol == "" {
		return path.Join(p.site, p.user, p.repo)
	}

	return fmt.Sprint(p.protocol, "://", path.Join(p.site, p.user, p.repo))
}

func describeGitHubRepo(uri string) (string, string, string, string, string, bool) {
	var (
		protocol = ""
		site     string
		user     string
		repo     string
		folder   = ""
	)

	startIndex := 0
	i := strings.Index(uri, "://")
	if i != -1 {
		protocol = uri[:i]
		startIndex = i + 3
	}

	bareURL := uri[startIndex:]
	bareParts := strings.Split(bareURL, "/")
	if bareParts[0] != "github.com" {
		return "", "", "", "", "", false
	} else if len(bareParts) < 3 {
		return "", "", "", "", "", false
	}

	site = bareParts[0]
	user = bareParts[1]
	repo = bareParts[2]

	if len(bareParts) > 3 {
		folder = strings.Join(bareParts[3:], "/")
	}

	return protocol, site, user, repo, folder, true
}

// DownloadFromGitHub takes a repo, and if it does not exists
// locally, it clones it to cache directory
func DownloadFromGitHub(pimpl remoteHandler, uri string, path string, debug bool) error {
	if ok := ShouldCloneProceed(pimpl, path); ok == false {
		return fmt.Errorf("repository already exists")
	}

	if err := pimpl.Download(uri, path, debug); err != nil {
		return err
	}

	return nil
}

// ShouldCloneProceed returns true if repository does not exist
// in specified folder
func ShouldCloneProceed(pimpl remoteHandler, path string) bool {
	if err := pimpl.Open(path); err != nil {
		return true
	}
	return false
}


type gitV4Handler struct{ remoteHandler }

func (g gitV4Handler) Download(uri string, dir string, debug bool) error {
	options := git.CloneOptions{
		URL: uri,
	}
	if debug {
		options.Progress = os.Stdout
	}
	_, err := git.PlainClone(dir, false, &options)
	return err
}

func (g gitV4Handler) Open(path string) error {
	_, err := git.PlainOpen(path)
	return err
}
