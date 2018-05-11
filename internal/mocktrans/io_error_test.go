package mocktrans

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	msg := "new error"
	err := NewError(msg, true, true)
	require.Implements(t, (*net.Error)(nil), err)
}

func TestIoError_Error(t *testing.T) {
	msg := "new error"
	err := NewError(msg, true, true)
	require.Equal(t, msg+" temp=true timeout=true", err.Error())
}

func TestIoError_Temporary(t *testing.T) {
	err := NewError("", true, false)
	require.True(t, err.Temporary())

	err = NewError("", false, false)
	require.False(t, err.Temporary())
}

func TestIoError_Timeout(t *testing.T) {
	err := NewError("", false, true)
	require.True(t, err.Timeout())

	err = NewError("", false, false)
	require.False(t, err.Timeout())
}
