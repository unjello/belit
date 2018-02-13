// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
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

		var err error

		compileCommand := []string{"/usr/bin/clang++", args[0], "-o", tempFile.Name()}
		if viper.GetBool("debug") {
			compileCommand = append(compileCommand, "-v")
		}

		if err = helpers.PrintCommand(compileCommand, viper.GetBool("debug")); err != nil {
			panic(err)
		}
		if err = helpers.PrintCommand([]string{tempFile.Name()}, true); err != nil {
			panic(err)
		}

		if err = helpers.RemoveFile(tempFile.Name()); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
