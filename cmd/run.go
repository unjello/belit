// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var AppFs = afero.NewOsFs()

// rootCmd represents the base command when called without any subcommands
var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Compile and run C/C++ program",
	Long:  `Compiles and runs C/C++ program`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}
		var (
			stdout io.ReadCloser
			stderr io.ReadCloser
			err    error
		)
		tempDir := afero.GetTempDir(AppFs, "belit")
		log.Info("Creating temporary directory: ", tempDir)
		tempFile, err2 := afero.TempFile(AppFs, tempDir, "belit-")
		if err2 != nil {
			panic(err2)
		}
		log.Info("Created temporary file: ", tempFile.Name())
		command := []string{"/usr/bin/clang++", args[0], "-o", tempFile.Name()}
		log.Info("Executing command: ", command)

		c := exec.Command(command[0], command[1:]...)

		if stderr, err = c.StderrPipe(); err != nil {
			panic(err)
		}
		if stdout, err = c.StdoutPipe(); err != nil {
			panic(err)
		}

		if err = c.Start(); err != nil {
			panic(err)
		}

		stderrOut, _ := ioutil.ReadAll(stderr)
		stdoutOut, _ := ioutil.ReadAll(stdout)

		if len(stderrOut) > 0 {
			fmt.Println(string(stderrOut))
		}
		if len(stdoutOut) > 0 {
			fmt.Println(string(stdoutOut))
		}

		if err = c.Wait(); err != nil {
			panic(err)
		}

		c = exec.Command(tempFile.Name())
		if stderr, err = c.StderrPipe(); err != nil {
			panic(err)
		}
		if stdout, err = c.StdoutPipe(); err != nil {
			panic(err)
		}

		if err = c.Start(); err != nil {
			panic(err)
		}

		stderrOut, _ = ioutil.ReadAll(stderr)
		stdoutOut, _ = ioutil.ReadAll(stdout)

		if len(stderrOut) > 0 {
			fmt.Println(string(stderrOut))
		}
		if len(stdoutOut) > 0 {
			fmt.Println(string(stdoutOut))
		}

		if err = c.Wait(); err != nil {
			panic(err)
		}

		err = os.Remove(tempFile.Name())
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
