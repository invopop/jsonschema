package jsonschema

import (
	"testing"
)

const (
	baseURL = "github.com/invopop/jsonschema"
)

// WantStruct is a test struct
// for testing reflection and
// ensuring that comment newlines
// are preserved when using
// jsonschema.NoSynopsis().
type WantStruct struct {
}

func TestDescription(t *testing.T) {
	r := new(Reflector)
	r.DoNotReference = true
	if err := r.AddGoComments(baseURL, "./"); err != nil {
		t.Fatal(err)
	}
	v := &WantStruct{}
	s := r.Reflect(v)

	want := `WantStruct is a test struct for testing reflection and ensuring that comment newlines are preserved when using jsonschema.NoSynopsis().`

	if got := s.Description; got != want {
		t.Errorf("s.Description =\n%v\nwant:\n%v", got, want)
	}
}

func TestDescription_NoSynopsis(t *testing.T) {
	r := new(Reflector)
	r.DoNotReference = true
	if err := r.AddGoComments(baseURL, "./", NoSynopsis()); err != nil {
		t.Fatal(err)
	}
	v := &WantStruct{}
	s := r.Reflect(v)

	want := `WantStruct is a test struct
for testing reflection and
ensuring that comment newlines
are preserved when using
jsonschema.NoSynopsis().`

	if got := s.Description; got != want {
		t.Errorf("s.Description =\n%v\nwant:\n%v", got, want)
	}
}
