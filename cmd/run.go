// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unjello/belit/config"
	"github.com/unjello/belit/helpers"
	src "github.com/unjello/belit/sources"
)

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
		baseDir := path.Join(viper.GetString("base"), "src")
		log := config.GetConfig().Log

		var includes []string
		for _, s := range sources {
			log.WithFields(logrus.Fields{
				"file":   fileName,
				"repo":   s.RepositoryPath,
				"header": s.HeaderName,
			}).Debug("Found header meta embedded in source code.")

			provider := src.GitProvider{}
			if provider.CanHandle(s.RepositoryPath) {
				inc := fmt.Sprintf("-I%s", provider.GetIncludePath(baseDir))
				includes = append(includes, inc)
			} else {
				panic("Cannot handle repo")
			}
		}

		compileCommand := []string{viper.GetString(meta.CompilerEnv), fileName, "-o", tempFile}
		compileOptionsStr := viper.GetString(meta.CompilerOptionsEnv)
		// `strings.Split` does return 1-element array if string is empty, but separator not.
		// need to test for string length first.
		buf, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		modeline, err := helpers.ExtractModelineOptions(string(buf))
		if err != nil {
			panic(err)
		}
		// give modeline a priority over options coming from env or config
		if len(modeline) > 0 {
			compileOptions := strings.Split(modeline[meta.CompilerOptionsEnv], " ")
			compileCommand = append(compileCommand, compileOptions...)
		} else if len(compileOptionsStr) > 0 {
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
