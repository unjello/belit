package config

import "github.com/sirupsen/logrus"

// Logger is an interface with richer methods set than StdLogger
type Logger logrus.FieldLogger

// RemoteRepository is an interface for downloadning different kinds of repositories
type RemoteRepository interface {
	Download(log Logger, url string, path string) error
}
