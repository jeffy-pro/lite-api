package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIErr_Error(t *testing.T) {
	err := NewAPIErr("some-code", "some-message")
	assert.Equal(t, "code: some-code | message: some-message", err.Error())
}
