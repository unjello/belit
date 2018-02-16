// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unjello/belit/helpers"
)

// rootCmd represents the base command when called without any subcommands
var getCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Get package",
	Long: `Download package from remote repository, and cache it locally.
Currently only github repositories are supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}
		path := args[0]
		err := helpers.FileExists(path)
		var repos []string
		if err != nil {
			repos = []string{path}
		} else {
			// TODO: Parse repos info from file
		}
		// TODO: do not bellyup when repo exists
		for _, repo := range repos {
			err := helpers.DownloadRemote("/Users/angelo/.belit", repo)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
