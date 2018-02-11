// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"os"

	"github.com/unjello/belit/helpers"

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

		tempDir := afero.GetTempDir(AppFs, "belit")
		log.Info("Creating temporary directory: ", tempDir)
		tempFile, err2 := afero.TempFile(AppFs, tempDir, "belit-")
		if err2 != nil {
			panic(err2)
		}
		log.Info("Created temporary file: ", tempFile.Name())

		compileStdOut, compileStdErr, err := helpers.RunCommand([]string{"/usr/bin/clang++", args[0], "-o", tempFile.Name()})

		if err != nil {
			if len(compileStdErr) > 0 {
				fmt.Println(compileStdErr)
			}

			panic(err)
		}

		if len(compileStdOut) > 0 {
			fmt.Println(compileStdOut)
		}

		runStdOut, _, runErr := helpers.RunCommand([]string{tempFile.Name()})
		if runErr != nil {
			panic(runErr)
		}

		if len(runStdOut) > 0 {
			fmt.Println(runStdOut)
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
