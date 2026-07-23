package jsonschema

import "testing"

func TestNilReflector(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	var r *Reflector
	s := r.Reflect(struct{ Name string }{})
	if s == nil {
		t.Fatal("want schema")
	}
}
