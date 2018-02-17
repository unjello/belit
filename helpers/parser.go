// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type SourceInfo struct {
	RepositoryPath string
	HeaderName     string
}

const matchRepoUrl = `^\s*?#include\s*?\/\*\s*?([-a-zA-Z0-9@:%_\+\.~#?&\/=]+)\s*?\*\/\s*[<"](.*?)[>"]`

func GetSources(file string) ([]SourceInfo, error) {
	var sources []SourceInfo

	l := log.WithFields(log.Fields{
		"file": file,
	})
	l.Debug("Open file")
	f, err := os.Open(file) // For read access.
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Could not open the file")
	}
	scanner := bufio.NewScanner(f)
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
			sources = append(sources, SourceInfo{match[1], match[2]})
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Could not close scanner")
	}

	return sources, nil
}
