package types

import "github.com/invopop/jsonschema/internal/testdata/types/deeper"

type Dummy struct {
	A     string
	Dummy deeper.Dummy
}
