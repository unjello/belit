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

func PrintCommand(command []string, isDebug bool) error {
	stdout, stderr, err := RunCommand(command)

	if isDebug {
		if len(stderr) > 0 {
			fmt.Println(stderr)
		}
		if len(stdout) > 0 {
			fmt.Println(stdout)
		}
	}

	if err != nil {
		if isDebug == false && len(stderr) > 0 {
			fmt.Println(stderr)
		}
		return err
	}

	return nil
}

func RunCommand(command []string) (string, string, error) {
	var (
		stdout io.ReadCloser
		stderr io.ReadCloser
		err    error
	)

	if len(command) < 1 {
		return "", "", fmt.Errorf(errorNoCommand)
	}
	log.WithFields(log.Fields{
		"cmd": strings.Join(command, " "),
	}).Info(infoExecutingCommand)

	c := exec.Command(command[0], command[1:]...)

	if stderr, err = c.StderrPipe(); err != nil {
		log.Panic(panicFailedToGetStderrPipe)
		return "", "", err
	}
	if stdout, err = c.StdoutPipe(); err != nil {
		log.Panic(panicFailedToGetStdoutPipe)
		return "", "", err
	}

	if err = c.Start(); err != nil {
		log.Panic(panicFailedToStartCommand)
		return "", "", err
	}

	stderrOut, _ := ioutil.ReadAll(stderr)
	stdoutOut, _ := ioutil.ReadAll(stdout)

	if err = c.Wait(); err != nil {
		log.Error(panicCommandFailedToRun)
		return string(stdoutOut), string(stderrOut), err
	}

	return string(stdoutOut), string(stderrOut), nil
}
