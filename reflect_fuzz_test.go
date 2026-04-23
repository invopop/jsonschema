package jsonschema

import (
	"reflect"
	"testing"
)

// FuzzReflectFromTypeExpandedStruct exercises ReflectFromType against a
// fixed pool of reflect.Types while fuzzing the Reflector flag combinations.
// It guards against regressions of a nil-pointer panic that occurred when
// ExpandedStruct=true was combined with a non-struct reflect.Type, since
// only struct types register themselves in the local definitions map.
func FuzzReflectFromTypeExpandedStruct(f *testing.F) {
	// Seed byte layout: data[0] selects a type variant, data[1] is a flags
	// word (bit 3 = ExpandedStruct). Seeds below set ExpandedStruct=true on
	// types that do NOT register in definitions, which previously triggered
	// the nil deref.
	f.Add([]byte{3, 8}) // map[string]any + ExpandedStruct=true
	f.Add([]byte{4, 8}) // []string       + ExpandedStruct=true
	f.Add([]byte{5, 8}) // any            + ExpandedStruct=true
	f.Add([]byte{7, 8}) // enum-tagged field + ExpandedStruct=true

	types := []reflect.Type{
		reflect.TypeOf(struct{}{}),
		reflect.TypeOf(struct{ Name string }{}),
		reflect.TypeOf(struct {
			A string `json:"a" jsonschema:"required,minLength=1,maxLength=100"`
			B int    `json:"b,omitempty"`
		}{}),
		reflect.TypeOf(map[string]any{}),
		reflect.TypeOf([]string{}),
		reflect.TypeOf((*any)(nil)).Elem(),
		reflect.TypeOf(struct {
			Self *struct{ X int } `json:"self,omitempty"`
		}{}),
		reflect.TypeOf(struct {
			Tags string `jsonschema:"enum=a,enum=b,enum=c"`
		}{}.Tags),
	}

	f.Fuzz(func(_ *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		idx := int(data[0]) % len(types)
		typ := types[idx]

		r := &Reflector{}
		if len(data) > 1 {
			flags := data[1]
			r.AllowAdditionalProperties = flags&1 != 0
			r.RequiredFromJSONSchemaTags = flags&2 != 0
			r.DoNotReference = flags&4 != 0
			r.ExpandedStruct = flags&8 != 0
		}

		// Must not panic regardless of type or flags combination.
		_ = r.ReflectFromType(typ)
	})
}
