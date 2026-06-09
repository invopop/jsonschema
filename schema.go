package jsonschema

import (
	"encoding/json"

	orderedmap "github.com/pb33f/ordered-map/v2"
)

// Version is the JSON Schema version.
var Version = "https://json-schema.org/draft/2020-12/schema"

// Schema represents a JSON Schema object type.
// RFC draft-bhutton-json-schema-00 section 4.3
type Schema struct {
	// RFC draft-bhutton-json-schema-00
	Version     string      `json:"$schema,omitempty" yaml:"$schema,omitempty"`         // section 8.1.1
	ID          ID          `json:"$id,omitempty" yaml:"$id,omitempty"`                 // section 8.2.1
	Anchor      string      `json:"$anchor,omitempty" yaml:"$anchor,omitempty"`         // section 8.2.2
	Ref         string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`               // section 8.2.3.1
	DynamicRef  string      `json:"$dynamicRef,omitempty" yaml:"$dynamicRef,omitempty"` // section 8.2.3.2
	Definitions Definitions `json:"$defs,omitempty" yaml:"$defs,omitempty"`             // section 8.2.4
	Comments    string      `json:"$comment,omitempty" yaml:"$comment,omitempty"`       // section 8.3
	// RFC draft-bhutton-json-schema-00 section 10.2.1 (Sub-schemas with logic)
	AllOf []*Schema `json:"allOf,omitempty" yaml:"allOf,omitempty"` // section 10.2.1.1
	AnyOf []*Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"` // section 10.2.1.2
	OneOf []*Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"` // section 10.2.1.3
	Not   *Schema   `json:"not,omitempty" yaml:"not,omitempty"`     // section 10.2.1.4
	// RFC draft-bhutton-json-schema-00 section 10.2.2 (Apply sub-schemas conditionally)
	If               *Schema            `json:"if,omitempty" yaml:"if,omitempty"`                             // section 10.2.2.1
	Then             *Schema            `json:"then,omitempty" yaml:"then,omitempty"`                         // section 10.2.2.2
	Else             *Schema            `json:"else,omitempty" yaml:"else,omitempty"`                         // section 10.2.2.3
	DependentSchemas map[string]*Schema `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty"` // section 10.2.2.4
	// RFC draft-bhutton-json-schema-00 section 10.3.1 (arrays)
	PrefixItems []*Schema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"` // section 10.3.1.1
	Items       *Schema   `json:"items,omitempty" yaml:"items,omitempty"`             // section 10.3.1.2  (replaces additionalItems)
	Contains    *Schema   `json:"contains,omitempty" yaml:"contains,omitempty"`       // section 10.3.1.3
	// RFC draft-bhutton-json-schema-00 section 10.3.2 (sub-schemas)
	Properties           *orderedmap.OrderedMap[string, *Schema] `json:"properties,omitempty" yaml:"properties,omitempty"`                     // section 10.3.2.1
	PatternProperties    map[string]*Schema                      `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`       // section 10.3.2.2
	AdditionalProperties *Schema                                 `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"` // section 10.3.2.3
	PropertyNames        *Schema                                 `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`               // section 10.3.2.4
	// RFC draft-bhutton-json-schema-validation-00, section 6
	Type              string              `json:"type,omitempty" yaml:"type,omitempty"`                           // section 6.1.1
	Enum              []any               `json:"enum,omitempty" yaml:"enum,omitempty"`                           // section 6.1.2
	Const             any                 `json:"const,omitempty" yaml:"const,omitempty"`                         // section 6.1.3
	MultipleOf        json.Number         `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`               // section 6.2.1
	Maximum           json.Number         `json:"maximum,omitempty" yaml:"maximum,omitempty"`                     // section 6.2.2
	ExclusiveMaximum  json.Number         `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`   // section 6.2.3
	Minimum           json.Number         `json:"minimum,omitempty" yaml:"minimum,omitempty"`                     // section 6.2.4
	ExclusiveMinimum  json.Number         `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`   // section 6.2.5
	MaxLength         *uint64             `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`                 // section 6.3.1
	MinLength         *uint64             `json:"minLength,omitempty" yaml:"minLength,omitempty"`                 // section 6.3.2
	Pattern           string              `json:"pattern,omitempty" yaml:"pattern,omitempty"`                     // section 6.3.3
	MaxItems          *uint64             `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`                   // section 6.4.1
	MinItems          *uint64             `json:"minItems,omitempty" yaml:"minItems,omitempty"`                   // section 6.4.2
	UniqueItems       bool                `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`             // section 6.4.3
	MaxContains       *uint64             `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`             // section 6.4.4
	MinContains       *uint64             `json:"minContains,omitempty" yaml:"minContains,omitempty"`             // section 6.4.5
	MaxProperties     *uint64             `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`         // section 6.5.1
	MinProperties     *uint64             `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`         // section 6.5.2
	Required          []string            `json:"required,omitempty" yaml:"required,omitempty"`                   // section 6.5.3
	DependentRequired map[string][]string `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"` // section 6.5.4
	// RFC draft-bhutton-json-schema-validation-00, section 7
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// RFC draft-bhutton-json-schema-validation-00, section 8
	ContentEncoding  string  `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`   // section 8.3
	ContentMediaType string  `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"` // section 8.4
	ContentSchema    *Schema `json:"contentSchema,omitempty" yaml:"contentSchema,omitempty"`       // section 8.5
	// RFC draft-bhutton-json-schema-validation-00, section 9
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`             // section 9.1
	Description string `json:"description,omitempty" yaml:"description,omitempty"` // section 9.1
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`         // section 9.2
	Deprecated  bool   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`   // section 9.3
	ReadOnly    bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`       // section 9.4
	WriteOnly   bool   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`     // section 9.4
	Examples    []any  `json:"examples,omitempty" yaml:"examples,omitempty"`       // section 9.5

	Extras map[string]any `json:"-" yaml:"-"` // Additional properties not defined in the schema

	// Special boolean representation of the Schema - section 4.3.2
	boolean *bool
}

var (
	// TrueSchema defines a schema with a true value
	TrueSchema = &Schema{boolean: &[]bool{true}[0]}
	// FalseSchema defines a schema with a false value
	FalseSchema = &Schema{boolean: &[]bool{false}[0]}
)

// Definitions hold schema definitions.
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
// RFC draft-wright-json-schema-validation-00, section 5.26
type Definitions map[string]*Schema
