package jsonschema

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iancoleman/orderedmap"

	"github.com/invopop/jsonschema/examples"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateFixtures = flag.Bool("update", false, "set to update fixtures")
var compareFixtures = flag.Bool("compare", false, "output failed fixtures with .out.json")

type GrandfatherType struct {
	FamilyName string `json:"family_name" jsonschema:"required"`
}

type SomeBaseType struct {
	ID               string `json:"id"` // to test composition override
	SomeBaseProperty int    `json:"some_base_property"`
	// The jsonschema required tag is nonsensical for private and ignored properties.
	// Their presence here tests that the fields *will not* be required in the output
	// schema, even if they are tagged required.
	somePrivateBaseProperty   string          `jsonschema:"required"`
	SomeIgnoredBaseProperty   string          `json:"-" jsonschema:"required"`
	SomeSchemaIgnoredProperty string          `jsonschema:"-,required"`
	Grandfather               GrandfatherType `json:"grand"`

	SomeUntaggedBaseProperty           bool `jsonschema:"required"`
	someUnexportedUntaggedBaseProperty bool
}

type MapType map[string]interface{}

type ArrayType []string

type nonExported struct {
	PublicNonExported  int
	privateNonExported int
}

type ProtoEnum int32

func (ProtoEnum) EnumDescriptor() ([]byte, []int) { return []byte(nil), []int{0} }

const (
	Unset ProtoEnum = iota
	Great
)

type TestUser struct {
	SomeBaseType
	nonExported
	MapType

	ID       int                    `json:"id" jsonschema:"required"`
	Name     string                 `json:"name" jsonschema:"required,minLength=1,maxLength=20,pattern=.*,description=this is a property,title=the name,example=joe,example=lucy,default=alex,readOnly=true"`
	Password string                 `json:"password" jsonschema:"writeOnly=true"`
	Friends  []int                  `json:"friends,omitempty" jsonschema_description:"list of IDs, omitted when empty"`
	Tags     map[string]string      `json:"tags,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`

	TestFlag       bool
	TestFlagFalse  bool `json:",omitempty" jsonschema:"default=false"`
	TestFlagTrue   bool `json:",omitempty" jsonschema:"default=true"`
	IgnoredCounter int  `json:"-"`

	// Tests for RFC draft-wright-json-schema-validation-00, section 7.3
	BirthDate time.Time `json:"birth_date,omitempty"`
	Website   url.URL   `json:"website,omitempty"`
	IPAddress net.IP    `json:"network_address,omitempty"`

	// Tests for RFC draft-wright-json-schema-hyperschema-00, section 4
	Photo  []byte `json:"photo,omitempty" jsonschema:"required"`
	Photo2 Bytes  `json:"photo2,omitempty" jsonschema:"required"`

	// Tests for jsonpb enum support
	Feeling ProtoEnum `json:"feeling,omitempty"`

	Age   int    `json:"age" jsonschema:"minimum=18,maximum=120,exclusiveMaximum=true,exclusiveMinimum=true"`
	Email string `json:"email" jsonschema:"format=email"`
	UUID  string `json:"uuid" jsonschema:"format=uuid"`

	// Test for "extras" support
	Baz       string `jsonschema_extras:"foo=bar,hello=world,foo=bar1"`
	BoolExtra string `json:"bool_extra,omitempty" jsonschema_extras:"isTrue=true,isFalse=false"`

	// Tests for simple enum tags
	Color      string  `json:"color" jsonschema:"enum=red,enum=green,enum=blue"`
	Rank       int     `json:"rank,omitempty" jsonschema:"enum=1,enum=2,enum=3"`
	Multiplier float64 `json:"mult,omitempty" jsonschema:"enum=1.0,enum=1.5,enum=2.0"`

	// Tests for enum tags on slices
	Roles      []string  `json:"roles" jsonschema:"enum=admin,enum=moderator,enum=user"`
	Priorities []int     `json:"priorities,omitempty" jsonschema:"enum=-1,enum=0,enum=1,enun=2"`
	Offsets    []float64 `json:"offsets,omitempty" jsonschema:"enum=1.570796,enum=3.141592,enum=6.283185"`

	// Test for raw JSON
	Anything interface{}     `json:"anything,omitempty"`
	Raw      json.RawMessage `json:"raw"`
}

type CustomTime time.Time

type CustomTypeField struct {
	CreatedAt CustomTime
}

type CustomTimeWithInterface time.Time

type CustomTypeFieldWithInterface struct {
	CreatedAt CustomTimeWithInterface
}

func (CustomTimeWithInterface) JSONSchema() *Schema {
	return &Schema{
		Type:   "string",
		Format: "date-time",
	}
}

type RootOneOf struct {
	Field1 string      `json:"field1" jsonschema:"oneof_required=group1"`
	Field2 string      `json:"field2" jsonschema:"oneof_required=group2"`
	Field3 interface{} `json:"field3" jsonschema:"oneof_type=string;array"`
	Field4 string      `json:"field4" jsonschema:"oneof_required=group1"`
	Field5 ChildOneOf  `json:"child"`
}

type ChildOneOf struct {
	Child1 string      `json:"child1" jsonschema:"oneof_required=group1"`
	Child2 string      `json:"child2" jsonschema:"oneof_required=group2"`
	Child3 interface{} `json:"child3" jsonschema:"oneof_required=group2,oneof_type=string;array"`
	Child4 string      `json:"child4" jsonschema:"oneof_required=group1"`
}

type Text string

type TextNamed string

type Outer struct {
	TextNamed
	Text `json:",omitempty"`
	Inner
}

type OuterNamed struct {
	Text  `json:"text,omitempty"`
	Inner `json:"inner"`
}

type Inner struct {
	Foo string `yaml:"foo"`
}

type MinValue struct {
	Value int `json:"value4" jsonschema_extras:"minimum=0"`
}
type Bytes []byte

type TestNullable struct {
	Child1 string `json:"child1" jsonschema:"nullable"`
}

type CompactDate struct {
	Year  int
	Month int
}

type UserWithAnchor struct {
	Name string `json:"name" jsonschema:"anchor=Name"`
}

func (CompactDate) JSONSchema() *Schema {
	return &Schema{
		Type:        "string",
		Title:       "Compact Date",
		Description: "Short date that only includes year and month",
		Pattern:     "^[0-9]{4}-[0-1][0-9]$",
	}
}

type TestDescriptionOverride struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Age        uint   `json:"age"`
	MiddleName string `json:"middle_name,omitempty"`
}

func (TestDescriptionOverride) GetFieldDocString(fieldName string) string {
	switch fieldName {
	case "FirstName":
		return "test2"
	case "LastName":
		return "test3"
	case "Age":
		return "test4"
	case "MiddleName":
		return "test5"
	default:
		return ""
	}
}

type LookupName struct {
	Given   string `json:"first"`
	Surname string `json:"surname"`
}

type LookupUser struct {
	Name  *LookupName `json:"name"`
	Alias string      `json:"alias,omitempty"`
}

type CustomSliceOuter struct {
	Slice CustomSliceType `json:"slice"`
}

type CustomSliceType []string

func (CustomSliceType) JSONSchema() *Schema {
	return &Schema{
		OneOf: []*Schema{{
			Type: "string",
		}, {
			Type: "array",
			Items: &Schema{
				Type: "string",
			},
		}},
	}
}

type CustomMapType map[string]string

func (CustomMapType) JSONSchema() *Schema {
	properties := orderedmap.New()
	properties.Set("key", &Schema{
		Type: "string",
	})
	properties.Set("value", &Schema{
		Type: "string",
	})
	return &Schema{
		Type: "array",
		Items: &Schema{
			Type:       "object",
			Properties: properties,
			Required:   []string{"key", "value"},
		},
	}
}

type CustomMapOuter struct {
	MyMap CustomMapType `json:"my_map"`
}

type PatternTest struct {
	WithPattern string `json:"with_pattern" jsonschema:"minLength=1,pattern=[0-9]{1\\,4},maxLength=50"`
}

type RecursiveExample struct {
	Text  string              `json:"text"`
	Child []*RecursiveExample `json:"children,omitempty"`
}

type KeyNamedNested struct {
	NestedNotRenamedProperty string
	NotRenamed               string
}

type KeyNamed struct {
	ThisWasLeftAsIs      string
	NotComingFromJSON    bool           `json:"coming_from_json_tag_not_renamed"`
	NestedNotRenamed     KeyNamedNested `json:"nested_not_renamed"`
	UnicodeShenanigans   string
	RenamedByComputation int `jsonschema_description:"Description was preserved"`
}

func TestReflector(t *testing.T) {
	r := new(Reflector)
	s := "http://example.com/schema"
	r.SetBaseSchemaID(s)
	assert.EqualValues(t, s, r.BaseSchemaID)
}

func TestReflectFromType(t *testing.T) {
	r := new(Reflector)
	tu := new(TestUser)
	typ := reflect.TypeOf(tu)

	s := r.ReflectFromType(typ)
	assert.EqualValues(t, "https://github.com/invopop/jsonschema/test-user", s.ID)

	x := struct {
		Test string
	}{
		Test: "foo",
	}
	typ = reflect.TypeOf(x)
	s = r.Reflect(typ)
	assert.Empty(t, s.ID)
}

func TestSchemaGeneration(t *testing.T) {
	tests := []struct {
		typ       interface{}
		reflector *Reflector
		fixture   string
	}{
		{&TestUser{}, &Reflector{}, "fixtures/test_user.json"},
		{&UserWithAnchor{}, &Reflector{}, "fixtures/user_with_anchor.json"},
		{&TestUser{}, &Reflector{AssignAnchor: true}, "fixtures/test_user_assign_anchor.json"},
		{&TestUser{}, &Reflector{AllowAdditionalProperties: true}, "fixtures/allow_additional_props.json"},
		{&TestUser{}, &Reflector{RequiredFromJSONSchemaTags: true}, "fixtures/required_from_jsontags.json"},
		{&TestUser{}, &Reflector{ExpandedStruct: true}, "fixtures/defaults_expanded_toplevel.json"},
		{&TestUser{}, &Reflector{IgnoredTypes: []interface{}{GrandfatherType{}}}, "fixtures/ignore_type.json"},
		{&TestUser{}, &Reflector{DoNotReference: true}, "fixtures/no_reference.json"},
		{&TestUser{}, &Reflector{DoNotReference: true, AssignAnchor: true}, "fixtures/no_reference_anchor.json"},
		{&RootOneOf{}, &Reflector{RequiredFromJSONSchemaTags: true}, "fixtures/oneof.json"},
		{&TestUser{}, &Reflector{ReferenceRoot: "#/components/schemas/"}, "fixtures/test_user_reference_root.json"},
		{&CustomTypeField{}, &Reflector{
			Mapper: func(i reflect.Type) *Schema {
				if i == reflect.TypeOf(CustomTime{}) {
					return &Schema{
						Type:   "string",
						Format: "date-time",
					}
				}
				return nil
			},
		}, "fixtures/custom_type.json"},
		{LookupUser{}, &Reflector{BaseSchemaID: "https://example.com/schemas"}, "fixtures/base_schema_id.json"},
		{LookupUser{}, &Reflector{
			Lookup: func(i reflect.Type) ID {
				switch i {
				case reflect.TypeOf(LookupUser{}):
					return ID("https://example.com/schemas/lookup-user")
				case reflect.TypeOf(LookupName{}):
					return ID("https://example.com/schemas/lookup-name")
				}
				return EmptyID
			},
		}, "fixtures/lookup.json"},
		{&LookupUser{}, &Reflector{
			BaseSchemaID:   "https://example.com/schemas",
			ExpandedStruct: true,
			AssignAnchor:   true,
			Lookup: func(i reflect.Type) ID {
				switch i {
				case reflect.TypeOf(LookupUser{}):
					return ID("https://example.com/schemas/lookup-user")
				case reflect.TypeOf(LookupName{}):
					return ID("https://example.com/schemas/lookup-name")
				}
				return EmptyID
			},
		}, "fixtures/lookup_expanded.json"},
		{&Outer{}, &Reflector{ExpandedStruct: true}, "fixtures/inlining_inheritance.json"},
		{&OuterNamed{}, &Reflector{ExpandedStruct: true}, "fixtures/inlining_embedded.json"},
		{&OuterNamed{}, &Reflector{ExpandedStruct: true, AssignAnchor: true}, "fixtures/inlining_embedded_anchored.json"},
		{&MinValue{}, &Reflector{}, "fixtures/schema_with_minimum.json"},
		{&TestNullable{}, &Reflector{}, "fixtures/nullable.json"},
		{&GrandfatherType{}, &Reflector{
			AdditionalFields: func(r reflect.Type) []reflect.StructField {
				return []reflect.StructField{
					{
						Name:      "Addr",
						Type:      reflect.TypeOf((*net.IP)(nil)).Elem(),
						Tag:       "json:\"ip_addr\"",
						Anonymous: false,
					},
				}
			},
		}, "fixtures/custom_additional.json"},
		{&TestDescriptionOverride{}, &Reflector{}, "fixtures/test_description_override.json"},
		{&CompactDate{}, &Reflector{}, "fixtures/compact_date.json"},
		{&CustomSliceOuter{}, &Reflector{}, "fixtures/custom_slice_type.json"},
		{&CustomMapOuter{}, &Reflector{}, "fixtures/custom_map_type.json"},
		{&CustomTypeFieldWithInterface{}, &Reflector{}, "fixtures/custom_type_with_interface.json"},
		{&PatternTest{}, &Reflector{}, "fixtures/commas_in_pattern.json"},
		{&examples.User{}, prepareCommentReflector(t), "fixtures/go_comments.json"},
		{&RecursiveExample{}, &Reflector{}, "fixtures/recursive.json"},
		{&KeyNamed{}, &Reflector{
			KeyNamer: func(s string) string {
				switch s {
				case "ThisWasLeftAsIs":
					fallthrough
				case "NotRenamed":
					fallthrough
				case "nested_not_renamed":
					return s
				case "coming_from_json_tag_not_renamed":
					return "coming_from_json_tag"
				case "NestedNotRenamed":
					return "nested-renamed"
				case "NestedNotRenamedProperty":
					return "nested-renamed-property"
				case "UnicodeShenanigans":
					return "✨unicode✨  s̸̥͝h̷̳͒e̴̜̽n̸̡̿a̷̘̔n̷̘͐i̶̫̐ǵ̶̯a̵̘͒n̷̮̾s̸̟̓"
				case "RenamedByComputation":
					return fmt.Sprintf("%.2f", float64(len(s))+1/137.0)
				}
				return "unknown case"
			},
		}, "fixtures/keynamed.json"},
		{MapType{}, &Reflector{}, "fixtures/map_type.json"},
		{ArrayType{}, &Reflector{}, "fixtures/array_type.json"},
	}

	for _, tt := range tests {
		name := strings.TrimSuffix(filepath.Base(tt.fixture), ".json")
		t.Run(name, func(t *testing.T) {
			compareSchemaOutput(t,
				tt.fixture, tt.reflector, tt.typ,
			)
		})
	}
}

func prepareCommentReflector(t *testing.T) *Reflector {
	t.Helper()
	r := new(Reflector)
	err := r.AddGoComments("github.com/invopop/jsonschema", "./examples")
	require.NoError(t, err, "did not expect error while adding comments")
	return r
}

func TestBaselineUnmarshal(t *testing.T) {
	r := &Reflector{}
	compareSchemaOutput(t, "fixtures/test_user.json", r, &TestUser{})
}

func compareSchemaOutput(t *testing.T, f string, r *Reflector, obj interface{}) {
	t.Helper()
	expectedJSON, err := ioutil.ReadFile(f)
	require.NoError(t, err)

	actualSchema := r.Reflect(obj)
	actualJSON, _ := json.MarshalIndent(actualSchema, "", "  ") //nolint:errchkjson

	if *updateFixtures {
		_ = ioutil.WriteFile(f, actualJSON, 0600)
	}

	if !assert.JSONEq(t, string(expectedJSON), string(actualJSON)) {
		if *compareFixtures {
			_ = ioutil.WriteFile(strings.TrimSuffix(f, ".json")+".out.json", actualJSON, 0600)
		}
	}
}

func TestSplitOnUnescapedCommas(t *testing.T) {
	tests := []struct {
		strToSplit string
		expected   []string
	}{
		{`Hello,this,is\,a\,string,haha`, []string{`Hello`, `this`, `is,a,string`, `haha`}},
		{`hello,no\\,split`, []string{`hello`, `no\,split`}},
		{`string without commas`, []string{`string without commas`}},
		{`ünicode,𐂄,Ж\,П,ᠳ`, []string{`ünicode`, `𐂄`, `Ж,П`, `ᠳ`}},
		{`empty,,tag`, []string{`empty`, ``, `tag`}},
	}

	for _, test := range tests {
		actual := splitOnUnescapedCommas(test.strToSplit)
		require.Equal(t, test.expected, actual)
	}
}

func TestArrayFormat(t *testing.T) {
	type URIArray struct {
		TestURIs []string `jsonschema:"type=array,format=uri"`
	}

	r := new(Reflector)
	schema := r.Reflect(&URIArray{})
	d := schema.Definitions["URIArray"]
	require.NotNil(t, d)
	props := d.Properties
	require.NotNil(t, props)
	i, found := props.Get("TestURIs")
	require.True(t, found)

	p := i.(*Schema)
	pt := p.Items.Format
	require.Equal(t, pt, "uri")
}
