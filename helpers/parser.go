// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"

)

type SourceInfo struct {
	RepositoryPath string
	HeaderName     string
}

var (
	ErrCouldNotOpenFile = errors.New("belit/parser: could not open file for reading")
	ErrCouldNotCloseScanner = errors.New("belit/parser: could not close scanner")
)

const matchRepoUrl = `^\s*?#include\s*?\/\*\s*?([-a-zA-Z0-9@:%_\+\.~#?&\/=]+)\s*?\*\/\s*[<"](.*?)[>"]`

func GetSourcesFromBuffer(buffer []byte) ([]SourceInfo, error) {
	var sources []SourceInfo

	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)

	re := regexp.MustCompile(matchRepoUrl)
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if match != nil {
			sources = append(sources, SourceInfo{match[1], match[2]})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, ErrCouldNotCloseScanner
	}

	return sources, nil
}

func GetSources(file string) ([]SourceInfo, error) {
	buffer, err := GetFileContents(file)
	if err != nil {
		return []SourceInfo{}, ErrCouldNotOpenFile
	}
	return GetSourcesFromBuffer(buffer)
}
