// andrzej lichnerowicz, unlicensed (~public domain)

// Package provides system and configuration abstraction

package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Config represents collection of system and configuration objects
type Config struct {
	Fs  afero.Fs
	Log *logrus.Logger
}

var config *Config
var once sync.Once

// GetConfig returns system configuration
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			Fs:  afero.NewOsFs(),
			Log: logrus.New(),
		}
		config.Log.Out = os.Stdout
	})
	return config
}
