package jsonschema

import (
	"net/url"
	"strings"
)

// JSONPointer escapes the path given to present a valid JSON pointer.
// This is required when the path has special characters like `/` or `~`
func JSONPointer(path string) string {
	path = strings.ReplaceAll(path, "~", "~0")
	path = strings.ReplaceAll(path, "/", "~1")
	return url.PathEscape(path)
}
