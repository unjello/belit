// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"bufio"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unjello/belit/helpers"
)

const matchRepoUrl = `^\s*?#include\s*?\/\*\s*?([-a-zA-Z0-9@:%_\+\.~#?&\/=]+)\s*?\*\/\s*[<"](.*?)[>"]`

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
			l := log.WithFields(log.Fields{
				"file": path,
			})
			l.Debug("Open file")
			file, err := os.Open(path) // For read access.
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			l2 := log.WithFields(log.Fields{
				"pattern": matchRepoUrl,
			})
			l2.Debug("Look for headers using regexp")
			re := regexp.MustCompile(matchRepoUrl)
			for scanner.Scan() {
				match := re.FindStringSubmatch(scanner.Text())
				if match != nil {
					l.WithFields(log.Fields{
						"match": match[1],
					}).Debug("Found repository")
					repos = append(repos, match[1])
				}
			}
			if err = scanner.Err(); err != nil {
				panic(err)
			}
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
