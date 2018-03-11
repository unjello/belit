// andrzej lichnerowicz, unlicensed (~public domain)

package sources

import (
	"os"

	"gopkg.in/src-d/go-git.v4"
)

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
