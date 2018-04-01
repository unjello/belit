// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
)

const (
	matchEmbeddedModeline = `\/(\/|\*) belit: ([^\n\*]+)`
	matchModelineOptions  = `([a-zA-Z]+)=("([^"]+)"|([\S]+))`
)

// ExtractModelineOptions searches throght the file in search of vim's
// [modeline](http://vimdoc.sourceforge.net/htmldoc/options.html#modeline)
// compatibile options
func ExtractModelineOptions(buffer string) (map[string]string, error) {
	ret := make(map[string]string)

	modeline, err := extractModeline(buffer)
	if err != nil {
		return map[string]string{}, err
	}
	reader := strings.NewReader(modeline)
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(matchModelineOptions)
	for scanner.Scan() {
		match := re.FindAllStringSubmatch(scanner.Text(), 10)
		if match != nil {
			for _, v := range match {
				if v[3] != "" {
					ret[v[1]] = v[3]
				} else {
					ret[v[1]] = v[4]
				}
			}
		}
	}
	return ret, nil
}

var (
	ErrorModelineNotFound = errors.New("modeline not found")
)

func extractModeline(buffer string) (string, error) {
	reader := strings.NewReader(buffer)
	scanner := bufio.NewScanner(reader)

	re := regexp.MustCompile(matchEmbeddedModeline)
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if match != nil {
			// regexp with \s at the end (\/(\/|\*) belit: ([^\n\*]+)\s)
			// works nice in regexp101, but is too non-greedy in Golang
			// therefore I allow a space, but trim result
			return strings.TrimSpace(match[2]), nil
		}
	}

	return "", ErrorModelineNotFound
}
