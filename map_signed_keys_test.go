package jsonschema

import (
	"encoding/json"
	"math"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignedMapSchemaAcceptsMarshalledKeys(t *testing.T) {
	value := map[int64]string{math.MinInt64: "negative", 0: "zero", 42: "positive"}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	var object map[string]string
	require.NoError(t, json.Unmarshal(data, &object))

	schema := Reflect(value)
	require.Equal(t, FalseSchema, schema.AdditionalProperties)
	require.Len(t, schema.PatternProperties, 1)
	for pattern, item := range schema.PatternProperties {
		require.Equal(t, "string", item.Type)
		matcher, err := regexp.Compile(pattern)
		require.NoError(t, err)
		for key := range object {
			require.True(t, matcher.MatchString(key), "schema rejects marshalled key %q", key)
		}
		for _, key := range []string{"word", "1.5", "--1", "1x"} {
			require.False(t, matcher.MatchString(key), "schema accepts non-integer key %q", key)
		}
	}
}
