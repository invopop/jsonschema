package jsonschema

import "testing"

func TestIDValidateTrimSpace(t *testing.T) {
	id := ID("  https://example.com/schema  ")
	if err := id.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := ID("   ").Validate(); err == nil {
		t.Fatal("expected empty id error")
	}
}
