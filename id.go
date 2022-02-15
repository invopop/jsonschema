package jsonschema

import "strings"

// ID represents a Schema ID type which should always be a URI.
// See draft-bhutton-json-schema-00 section 8.2.1
type ID string

// Anchor either adds or replaces the anchor part of the schema URI.
func (id ID) Anchor(name string) ID {
	b := id.Base()
	return ID(b.String() + "#" + name)
}

// Def adds or replaces a definition identifier.
func (id ID) Def(name string) ID {
	b := id.Base()
	return ID(b.String() + "#/$defs/" + name)
}

// Add appends the provided path to the id, and removes any
// anchor data that might be there.
func (id ID) Add(path string) ID {
	b := id.Base()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return ID(b.String() + path)
}

// Base removes any anchor information from the schema
func (id ID) Base() ID {
	s := id.String()
	i := strings.LastIndex(s, "#")
	if i != -1 {
		s = s[0:i]
	}
	s = strings.TrimRight(s, "/")
	return ID(s)
}

// String provides string version of ID
func (id ID) String() string {
	return string(id)
}
