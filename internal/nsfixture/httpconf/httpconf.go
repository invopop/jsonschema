// Package httpconf provides a fixture type used to test AddPackageNamespaces:
// a Config struct that collides in name with tcpconf.Config.
package httpconf

// Config holds HTTP-specific configuration.
type Config struct {
	URL    string `json:"url" jsonschema:"required,format=uri"`
	Method string `json:"method" jsonschema:"required"`
}
