# Go JSON Schema Reflection

[![CI](https://github.com/invopop/jsonschema/actions/workflows/ci.yml/badge.svg)](https://github.com/invopop/jsonschema/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/jsonschema)](https://goreportcard.com/report/github.com/invopop/jsonschema)
[![GoDoc](https://godoc.org/github.com/invopop/jsonschema?status.svg)](https://godoc.org/github.com/invopop/jsonschema)

This package can be used to generate [JSON Schemas](http://json-schema.org/latest/json-schema-validation.html) from Go types through reflection.

- Supports arbitrarily complex types, including `interface{}`, maps, slices, etc.
- Supports json-schema features such as minLength, maxLength, pattern, format, etc.
- Supports simple string and numeric enums.
- Supports custom property fields via the `jsonschema_extras` struct tag.

This repository is a fork of the original [jsonschema](https://github.com/alecthomas/jsonschema) by [@alecthomas](https://github.com/alecthomas). At [Invopop](https://invopop.com) we use jsonschema as a cornerstone in our [GOBL library](https://github.com/invopop/gobl), and wanted to be able to continue building and adding features without taking up Alec's time.

## Example

The following Go type:

```go
type TestUser struct {
  ID            int                    `json:"id"`
  Name          string                 `json:"name" jsonschema:"title=the name,description=The name of a friend,example=joe,example=lucy,default=alex"`
  Friends       []int                  `json:"friends,omitempty" jsonschema_description:"The list of IDs, omitted when empty"`
  Tags          map[string]interface{} `json:"tags,omitempty" jsonschema_extras:"a=b,foo=bar,foo=bar1"`
  BirthDate     time.Time              `json:"birth_date,omitempty" jsonschema:"oneof_required=date"`
  YearOfBirth   string                 `json:"year_of_birth,omitempty" jsonschema:"oneof_required=year"`
  Metadata      interface{}            `json:"metadata,omitempty" jsonschema:"oneof_type=string;array"`
  FavColor      string                 `json:"fav_color,omitempty" jsonschema:"enum=red,enum=green,enum=blue"`
}
```

Results in following JSON Schema:

```go
jsonschema.Reflect(&TestUser{})
```

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$ref": "#/definitions/TestUser",
  "definitions": {
    "TestUser": {
      "type": "object",
      "properties": {
        "metadata": {
          "oneOf": [
            {
              "type": "string"
            },
            {
              "type": "array"
            }
          ]
        },
        "birth_date": {
          "type": "string",
          "format": "date-time"
        },
        "friends": {
          "type": "array",
          "items": {
            "type": "integer"
          },
          "description": "The list of IDs, omitted when empty"
        },
        "id": {
          "type": "integer"
        },
        "name": {
          "type": "string",
          "title": "the name",
          "description": "The name of a friend",
          "default": "alex",
          "examples": ["joe", "lucy"]
        },
        "tags": {
          "type": "object",
          "patternProperties": {
            ".*": {
              "additionalProperties": true
            }
          },
          "a": "b",
          "foo": ["bar", "bar1"]
        },
        "fav_color": {
          "type": "string",
          "enum": ["red", "green", "blue"]
        }
      },
      "additionalProperties": false,
      "required": ["id", "name"],
      "oneOf": [
        {
          "required": ["birth_date"],
          "title": "date"
        },
        {
          "required": ["year_of_birth"],
          "title": "year"
        }
      ]
    }
  }
}
```

## Configurable behaviour

The behaviour of the schema generator can be altered with parameters when a `jsonschema.Reflector`
instance is created.

### ExpandedStruct

If set to `true`, makes the top level struct not to reference itself in the definitions. But type passed should be a struct type.

eg.

```go
type GrandfatherType struct {
	FamilyName string `json:"family_name" jsonschema:"required"`
}

type SomeBaseType struct {
	SomeBaseProperty int `json:"some_base_property"`
	// The jsonschema required tag is nonsensical for private and ignored properties.
	// Their presence here tests that the fields *will not* be required in the output
	// schema, even if they are tagged required.
	somePrivateBaseProperty            string `json:"i_am_private" jsonschema:"required"`
	SomeIgnoredBaseProperty            string `json:"-" jsonschema:"required"`
	SomeSchemaIgnoredProperty          string `jsonschema:"-,required"`
	SomeUntaggedBaseProperty           bool   `jsonschema:"required"`
	someUnexportedUntaggedBaseProperty bool
	Grandfather                        GrandfatherType `json:"grand"`
}
```

will output:

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "required": ["some_base_property", "grand", "SomeUntaggedBaseProperty"],
  "properties": {
    "SomeUntaggedBaseProperty": {
      "type": "boolean"
    },
    "grand": {
      "$schema": "http://json-schema.org/draft-04/schema#",
      "$ref": "#/definitions/GrandfatherType"
    },
    "some_base_property": {
      "type": "integer"
    }
  },
  "type": "object",
  "definitions": {
    "GrandfatherType": {
      "required": ["family_name"],
      "properties": {
        "family_name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  }
}
```

### PreferYAMLSchema

JSON schemas can also be used to validate YAML, however YAML frequently uses
different identifiers to JSON indicated by the `yaml:` tag. The `Reflector` will
by default prefer `json:` tags over `yaml:` tags (and only use the latter if the
former are not present). This behavior can be changed via the `PreferYAMLSchema`
flag, that will switch this behavior: `yaml:` tags will be preferred over
`json:` tags.

With `PreferYAMLSchema: true`, the following struct:

```go
type Person struct {
	FirstName string `json:"FirstName" yaml:"first_name"`
}
```

would result in this schema:

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$ref": "#/definitions/TestYamlAndJson",
  "definitions": {
    "Person": {
      "required": ["first_name"],
      "properties": {
        "first_name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  }
}
```

whereas without the flag one obtains:

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$ref": "#/definitions/TestYamlAndJson",
  "definitions": {
    "Person": {
      "required": ["FirstName"],
      "properties": {
        "first_name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  }
}
```

### Using Go Comments

Writing a good schema with descriptions inside tags can become cumbersome and tedious, especially if you already have some Go comments around your types and field definitions. If you'd like to take advantage of these existing comments, you can use the `AddGoComments(base, path string)` method that forms part of the reflector to parse your go files and automatically generate a dictionary of Go import paths, types, and fields, to individual comments. These will then be used automatically as description fields, and can be overridden with a manual definition if needed.

Take a simplified example of a User struct which for the sake of simplicity we assume is defined inside this package:

```go
package main

// User is used as a base to provide tests for comments.
type User struct {
	// Unique sequential identifier.
	ID int `json:"id" jsonschema:"required"`
	// Name of the user
	Name string `json:"name"`
}
```

To get the comments provided into your JSON schema, use a regular `Reflector` and add the go code using an import module URL and path. Fully qualified go module paths cannot be determined reliably by the `go/parser` library, so we need to introduce this manually:

```go
r := new(Reflector)
if err := r.AddGoComments("github.com/invopop/jsonschema", "./"); err != nil {
  // deal with error
}
s := r.Reflect(&User{})
// output
```

Expect the results to be similar to:

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$ref": "#/definitions/User",
  "definitions": {
    "User": {
      "required": ["id"],
      "properties": {
        "id": {
          "type": "integer",
          "description": "Unique sequential identifier."
        },
        "name": {
          "type": "string",
          "description": "Name of the user"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "description": "User is used as a base to provide tests for comments."
    }
  }
}
```

### Custom Type Definitions

Sometimes it can be useful to have custom JSON Marshal and Unmarshal methods in your structs that automatically convert for example a string into an object.

To override auto-generating an object type for your type, implement the `JSONSchemaType() *Type` method and whatever is defined will be provided in the schema definitions.

Take the following simplified example of a `CompactDate` that only includes the Year and Month:

```go
type CompactDate struct {
	Year  int
	Month int
}

func (d *CompactDate) UnmarshalJSON(data []byte) error {
  if len(data) != 9 {
    return errors.New("invalid compact date length")
  }
  var err error
  d.Year, err = strconv.Atoi(string(data[1:5]))
  if err != nil {
    return err
  }
  d.Month, err = strconv.Atoi(string(data[7:8]))
  if err != nil {
    return err
  }
  return nil
}

func (d *CompactDate) MarshalJSON() ([]byte, error) {
  buf := new(bytes.Buffer)
  buf.WriteByte('"')
  buf.WriteString(fmt.Sprintf("%d-%02d", d.Year, d.Month))
  buf.WriteByte('"')
  return buf.Bytes(), nil
}

func (CompactDate) JSONSchemaType() *Type {
	return &Type{
		Type:        "string",
		Title:       "Compact Date",
		Description: "Short date that only includes year and month",
		Pattern:     "^[0-9]{4}-[0-1][0-9]$",
	}
}
```

The resulting schema generated for this struct would look like:

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$ref": "#/definitions/CompactDate",
  "definitions": {
    "CompactDate": {
      "pattern": "^[0-9]{4}-[0-1][0-9]$",
      "type": "string",
      "title": "Compact Date",
      "description": "Short date that only includes year and month"
    }
  }
}
```
