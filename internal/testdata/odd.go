package testdata

import (
	"github.com/invopop/jsonschema/internal/testdata/types"
	"github.com/invopop/jsonschema/internal/testdata/types/deeper"
)

type (
	Odd struct {
		Dummy1 types.Dummy  `json:"dummy1"`
		Dummy2 deeper.Dummy `json:"dummy2"`
		Dummy3 Dummy        `json:"dummy3"`
	}
	Dummy struct {
		B int
	}
)
