package jsonschema

import (
	"net/url"
	"strings"
)

func JSONPointer(path string) string {
	path = strings.ReplaceAll(path, "~", "~0")
	path = strings.ReplaceAll(path, "/", "~1")
	return url.PathEscape(path)
}
