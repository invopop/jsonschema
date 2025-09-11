package jsonschema

import (
	orderedmap "github.com/pb33f/ordered-map/v2"
)

// NewProperties is a helper method to instantiate a new properties ordered
// map.
func NewProperties() *orderedmap.OrderedMap[string, *Schema] {
	return orderedmap.New[string, *Schema]()
}
