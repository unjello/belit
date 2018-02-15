// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"os"

	"github.com/unjello/belit/helpers"

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

		var (
			err      error
			tempFile string
		)

		if tempFile, err = helpers.GetTempFile(); err != nil {
			panic(err)
		}

		// TODO: add CCFLAGS
		// TODO: add CXXFLAGS
		fileName := args[0]
		meta, err := helpers.GetCompilerOptions(fileName)
		if err != nil {
			panic(err)
		}
		compileCommand := []string{viper.GetString(meta.CompilerEnv), fileName, "-o", tempFile}
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
