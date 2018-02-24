// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unjello/belit/helpers"
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
			err      error
			tempFile string
		)

		if tempFile, err = helpers.GetTempFile(); err != nil {
			panic(err)
		}

		fileName := args[0]
		meta, err := helpers.GetCompilerOptions(fileName)
		if err != nil {
			panic(err)
		}

		sources, err := helpers.GetSources(fileName)
		if err != nil {
			panic(err)
		}
		// TODO: Refactor this into config
		baseDir := "/Users/angelo/.belit/src"
		var includes []string
		for _, s := range sources {
			log.WithFields(log.Fields{
				"file":   fileName,
				"repo":   s.RepositoryPath,
				"header": s.HeaderName,
			}).Debug("Found header meta embedded in source code.")

			inc := fmt.Sprintf("-I%s", path.Join(baseDir, s.RepositoryPath))
			includes = append(includes, inc)
			log.WithFields(log.Fields{
				"file":    fileName,
				"include": inc,
			}).Info("Adding compiler options")
		}

		compileCommand := []string{viper.GetString(meta.CompilerEnv), fileName, "-std=c++11", "-o", tempFile}
		compileOptionsStr := viper.GetString(meta.CompilerOptionsEnv)
		// `strings.Split` does return 1-element array if string is empty, but separator not.
		// need to test for string length first.
		if len(compileOptionsStr) > 0 {
			compileOptions := strings.Split(compileOptionsStr, " ")
			compileCommand = append(compileCommand, compileOptions...)
		}
		if len(includes) > 0 {
			compileCommand = append(compileCommand, includes...)
		}
		if viper.GetBool("debug") {
			compileCommand = append(compileCommand, "-v")
		}

		if err = helpers.PrintCommand(compileCommand, viper.GetBool("debug")); err != nil {
			panic(err)
		}
		if err = helpers.PrintCommand([]string{tempFile}, true); err != nil {
			panic(err)
		}

		if err = helpers.RemoveFile(tempFile); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
