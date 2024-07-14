package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIErr_Error(t *testing.T) {
	err := NewAPIErr("some-code", "some-message")
	require.Equal(t, "code: some-code | message: some-message", err.Error())
}
