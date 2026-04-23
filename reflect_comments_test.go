package jsonschema

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/invopop/jsonschema/examples"
)

func TestCommentsSchemaGeneration(t *testing.T) {
	tests := []struct {
		typ       any
		reflector *Reflector
		fixture   string
	}{
		{&examples.User{}, prepareCommentReflector(t), "fixtures/go_comments.json"},
		{&examples.User{}, prepareCommentReflector(t, WithFullComment()), "fixtures/go_comments_full.json"},
		{&examples.User{}, prepareCustomCommentReflector(t), "fixtures/custom_comments.json"},
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

func TestAddGoCommentsSkipsUnexportedTypes(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

// exportedDoc documents the Exported type.
type Exported struct {
	// Field is a field comment.
	Field string
}

// unexportedDoc documents the unexported type.
type unexported struct {
	Field string
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sample.go"), []byte(src), 0o600))

	r := new(Reflector)
	require.NoError(t, r.AddGoComments("example.com/sample", dir))

	var exportedFound bool
	for k := range r.CommentMap {
		require.NotContains(t, k, "unexported", "unexported type must not appear in comment map")
		if strings.HasSuffix(k, ".Exported") {
			exportedFound = true
		}
	}
	require.True(t, exportedFound, "exported type should have been picked up")
}

func prepareCustomCommentReflector(t *testing.T) *Reflector {
	t.Helper()
	r := new(Reflector)
	r.LookupComment = func(t reflect.Type, f string) string {
		if t != reflect.TypeOf(examples.User{}) {
			// To test the interaction between a custom LookupComment function and the
			// AddGoComments function, we only override comments for the User type.
			return ""
		}
		if f == "" {
			return fmt.Sprintf("Go type %s, defined in package %s.", t.Name(), t.PkgPath())
		}
		return fmt.Sprintf("Field %s of Go type %s.%s.", f, t.PkgPath(), t.Name())
	}
	// Also add the Go comments.
	err := r.AddGoComments("github.com/invopop/jsonschema", "./examples")
	require.NoError(t, err, "did not expect error while adding comments")
	return r
}
