// Package tcpconf provides a fixture type used to test AddPackageNamespaces:
// a Config struct that collides in name with httpconf.Config.
package tcpconf

// Config holds TCP-specific configuration.
type Config struct {
	Host string `json:"host" jsonschema:"required,format=hostname"`
	Port int    `json:"port" jsonschema:"required,minimum=1,maximum=65535"`
}
