package jsonschema_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	base := "https://invopop.com/schema"
	id := jsonschema.ID(base)

	assert.Equal(t, base, id.String())

	id = id.Add("user")
	assert.EqualValues(t, base+"/user", id)

	id = id.Anchor("Name")
	assert.EqualValues(t, base+"/user#Name", id)

	id = id.Anchor("Title")
	assert.EqualValues(t, base+"/user#Title", id)

	id = id.Def("Name")
	assert.EqualValues(t, base+"/user#/$defs/Name", id)
}
