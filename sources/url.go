package sources

import (
	"fmt"
	"strings"
)

func ensureHTTPS(url string) string {
	i := strings.Index(url, "://")
	if i == -1 {
		return fmt.Sprint("https://", url)
	} else if url[:i] != "https" {
		return fmt.Sprint("https://", url[i+3:])
	}
	return url
}
