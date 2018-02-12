// andrzej lichnerowicz, unlicensed (~public domain)

package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/unjello/belit/cmd"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	cmd.Execute()
}
