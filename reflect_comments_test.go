package jsonschema

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/jsonschema/examples"
	"github.com/stretchr/testify/require"
)

func TestCommentsSchemaGeneration(t *testing.T) {
	tests := []struct {
		typ       any
		reflector *Reflector
		fixture   string
	}{
		{&examples.User{}, prepareCommentReflector(t), "fixtures/go_comments.json"},
		{&examples.User{}, prepareCommentReflector(t, WithFullComment()), "fixtures/go_comments_full.json"},
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

func prepareCommentReflector(t *testing.T, opts ...CommentOption) *Reflector {
	t.Helper()
	r := new(Reflector)
	err := r.AddGoComments("github.com/invopop/jsonschema", "./examples", opts...)
	require.NoError(t, err, "did not expect error while adding comments")
	return r
}
