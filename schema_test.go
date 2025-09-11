package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaMerge(t *testing.T) {
	t.Parallel()

	r := Reflector{
		ExpandedStruct: true,
	}

	combined, err := r.Reflect(&TestUser{}).
		MergeSchemas(r.Reflect(&Inner{}), r.Reflect(&RecursiveExample{}))

	require.NoError(t, err)
	require.NotNil(t, combined)

	// recursive definition must be added manually when using ExpandedStruct
	combined.AddDefinition("RecursiveExample", r.Reflect(&RecursiveExample{}))

	type combinedStruct struct {
		TestUser         `json:",inline"`
		Inner            `json:",inline"`
		RecursiveExample `json:",inline"`
	}

	expected := r.Reflect(&combinedStruct{})

	expected.ID = combined.ID // IDs are expected to differ, everything else must be identical

	expectedJSON, err := expected.MarshalJSON()
	require.NoError(t, err)

	combinedJSON, err := combined.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, string(expectedJSON), string(combinedJSON))
}
