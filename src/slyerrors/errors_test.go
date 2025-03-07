package slyerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	val := NewValidation("startCode1")
	val.ValidateNotEmpty("empty", "")

	val2 := NewValidation("startCode2")
	val2.ValidateNotEmpty("empty2", "")

	val.Merge(val2.Error(), "subfield")

	errstr := val.Error().Error()

	assert.Equal(t, "error: kind Validation (4), code: stringEmpty, message: invalid request, invalid fields: [empty: code: stringEmpty | msg: ]|[subfield.empty2: code: stringEmpty | msg: ]", errstr)
}
