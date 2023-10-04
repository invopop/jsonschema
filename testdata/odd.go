package testdata

import (
	"github.com/invopop/jsonschema/testdata/pkg/types"
	types2 "github.com/invopop/jsonschema/testdata/pkg2/types"
)

type Odd struct {
	Dummy1 types.Dummy  `json:"dummy1"`
	Dummy2 types2.Dummy `json:"dummy2"`
}
