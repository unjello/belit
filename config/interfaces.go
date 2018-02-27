package config

import "github.com/sirupsen/logrus"

// Logger is an interface with richer methods set than StdLogger
type Logger logrus.FieldLogger
