// andrzej lichnerowicz, unlicensed (~public domain)
package helpers

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
)

// ExtractModelineOptions searches throght the file in search of vim's
// [modeline](http://vimdoc.sourceforge.net/htmldoc/options.html#modeline)
// compatibile options
func ExtractModelineOptions(buffer string) (map[string]string, error) {
	ret := make(map[string]string)
	return ret, nil
}

var (
	ErrorModelineNotFound = errors.New("modeline not found")
)

const (
	matchEmbeddedModeline = `\/(\/|\*) belit: ([^\n\*]+)`
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
