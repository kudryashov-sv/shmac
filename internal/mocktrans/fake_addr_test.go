package mocktrans

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFakeAddr(t *testing.T) {
	addr := newFakeAddr()
	require.Implements(t, (*net.Addr)(nil), addr)
}

func TestFakeAddr_Network(t *testing.T) {
	addr := newFakeAddr()
	require.Equal(t, "fake", addr.Network())
}

func TestFakeAddr_String(t *testing.T) {
	addr := newFakeAddr()
	require.Equal(t, "fake:address", addr.String())
}
