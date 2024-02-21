package testdata

import (
	"github.com/invopop/jsonschema/internal/testdata/types"
	"github.com/invopop/jsonschema/internal/testdata/types/deeper"
)

type (
	Odd struct {
		Dummy1  types.Dummy
		Dummy1a types.Dummy

		Dummy2  deeper.Dummy
		Dummy2a deeper.Dummy

		Dummy3  Dummy
		Dummy3a Dummy
	}
	Dummy struct {
		B int
	}
)
