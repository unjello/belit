// andrzej lichnerowicz, unlicensed (~public domain)

package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	errorNoCommand             = "No command to execute"
	infoExecutingCommand       = "Executing command: "
	panicFailedToGetStderrPipe = "Failed to get stderr pipe"
	panicFailedToGetStdoutPipe = "Failed to get stdout pipe"
	panicFailedToStartCommand  = "Failed to start command"
	panicCommandFailedToRun    = "Command failed to run"
)

func RunCommand(command []string) (string, string, error) {
	var (
		stdout io.ReadCloser
		stderr io.ReadCloser
		err    error
	)

	if len(command) < 1 {
		return "", "", fmt.Errorf(errorNoCommand)
	}

	l := log.WithFields(log.Fields{
		"cmd": command,
	})
	log.Info(infoExecutingCommand, strings.Join(command, " "))

	c := exec.Command(command[0], command[1:]...)

	if stderr, err = c.StderrPipe(); err != nil {
		l.Panic(panicFailedToGetStderrPipe)
		return "", "", err
	}
	if stdout, err = c.StdoutPipe(); err != nil {
		l.Panic(panicFailedToGetStdoutPipe)
		return "", "", err
	}

	if err = c.Start(); err != nil {
		l.Panic(panicFailedToStartCommand)
		return "", "", err
	}

	stderrOut, _ := ioutil.ReadAll(stderr)
	stdoutOut, _ := ioutil.ReadAll(stdout)

	if err = c.Wait(); err != nil {
		l.Error(panicCommandFailedToRun)
		return string(stdoutOut), string(stderrOut), err
	}

	return string(stdoutOut), string(stderrOut), nil
}
