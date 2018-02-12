// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"os"

	"github.com/unjello/belit/helpers"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		tempDir := afero.GetTempDir(AppFs, "belit")
		log.WithFields(log.Fields{
			"folder": tempDir,
		}).Debug("Created temporary folder")

		tempFile, err2 := afero.TempFile(AppFs, tempDir, "belit-")
		if err2 != nil {
			panic(err2)
		}
		log.WithFields(log.Fields{
			"file": tempFile.Name(),
		}).Debug("Created temporary file")

		compileCommand := []string{"/usr/bin/clang++", args[0], "-o", tempFile.Name()}
		if viper.GetBool("debug") {
			compileCommand = append(compileCommand, "-v")
		}

		var (
			stdout string
			stderr string
			err    error
		)
		stdout, stderr, err = helpers.RunCommand(compileCommand)

		if viper.GetBool("debug") {
			if len(stderr) > 0 {
				fmt.Println(stderr)
			}
			if len(stdout) > 0 {
				fmt.Println(stdout)
			}
		}

		if err != nil {
			if viper.GetBool("debug") == false {
				if len(stderr) > 0 {
					fmt.Println(stderr)
				}
			}

			panic(err)
		}

		stdout, stderr, err = helpers.RunCommand([]string{tempFile.Name()})
		if err != nil {
			panic(err)
		}

		if len(stdout) > 0 {
			fmt.Println(stdout)
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
