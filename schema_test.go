package jsonschema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ethExample = `
{
    "BlobsBundleV2": {
        "properties": {
            "blobs": {
                "items": {
                    "$ref": "#/components/schemas/bytes"
                },
                "title": "Blobs",
                "type": "array"
            },
            "commitments": {
                "items": {
                    "$ref": "#/components/schemas/bytes48"
                },
                "title": "Commitments",
                "type": "array"
            },
            "proofs": {
                "items": {
                    "$ref": "#/components/schemas/bytes48"
                },
                "title": "Proofs",
                "type": "array"
            }
        },
        "required": [
            "commitments",
            "proofs",
            "blobs"
        ],
        "title": "Blobs bundle object V2",
        "type": "object"
    },
    "ExecutionPayloadBodyV1": {
        "properties": {
            "transactions": {
                "$ref": "#/components/schemas/ExecutionPayloadV1/properties/transactions"
            },
            "withdrawals": {
                "items": {
                    "$ref": "#/components/schemas/WithdrawalV1"
                },
                "title": "Withdrawals",
                "type": [
                    "array",
                    "null"
                ]
            }
        },
        "required": [
            "transactions"
        ],
        "title": "Execution payload body object V1",
        "type": "object"
    }
}
`

func TestType_MarshalJSON(t *testing.T) {
	t.Parallel()

	filename := filepath.Join("fixtures", "multivalued_type.json")

	expectedJSON, err := os.ReadFile(filename)
	require.NoError(t, err)

	props := NewProperties()
	props.Set("other", &Schema{
		Type: TypeString,
	})
	props.Set("withdrawals", &Schema{
		Type: NewMultivaluedType(TypeArray, TypeNull),
	})

	schema := &Schema{
		Description: "test of Schema with multivalued property Type",
		Properties:  props,
		Required:    []string{"withdrawals"},
	}

	actualJSON, err := json.MarshalIndent(schema, "", "  ")
	require.NoError(t, err)

	if *updateFixtures {
		_ = os.WriteFile(filename, actualJSON, 0600)
	}

	if !assert.JSONEq(t, string(expectedJSON), string(actualJSON)) {
		if *compareFixtures {
			_ = os.WriteFile(strings.TrimSuffix(filename, ".json")+".out.json", actualJSON, 0600)
		}
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var defs Definitions
	require.NoError(t, json.Unmarshal([]byte(ethExample), &defs))

	withdrawals, ok := defs["ExecutionPayloadBodyV1"].Properties.Get("withdrawals")
	require.True(t, ok)
	assert.True(t, withdrawals.Type.IsMultivalued())
	assert.Equal(t, NewMultivaluedType(TypeArray, TypeNull), withdrawals.Type)

	proofs, ok := defs["BlobsBundleV2"].Properties.Get("proofs")
	require.True(t, ok)
	assert.False(t, proofs.Type.IsMultivalued())
	assert.Equal(t, TypeArray, proofs.Type)
}
