package mocktrans

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTransport(t *testing.T) {
	tr, err := NewTransport()
	require.NoError(t, err)
	require.NotNil(t, tr)
	require.Implements(t, (*net.Listener)(nil), tr)
}

func TestTransport_Addr(t *testing.T) {
	tr, err := NewTransport()
	require.NoError(t, err)

	addr := tr.Addr()
	require.NotNil(t, addr)
}

func TestTransport_Close(t *testing.T) {
	tr, err := NewTransport()
	require.NoError(t, err)

	err = tr.Close()
	require.NoError(t, err)
}

func TestTransport_DialAccept(t *testing.T) {
	type acceptRes struct {
		cn  net.Conn
		err error
	}

	tr, err := NewTransport()
	require.NoError(t, err)

	accCh := make(chan acceptRes)
	go func() {
		res := acceptRes{}
		res.cn, res.err = tr.Accept()
		accCh <- res
	}()

	cn, err := tr.Dial()
	require.NoError(t, err)
	require.NotEmpty(t, cn)

	err = cn.Close()
	require.NoError(t, err)

	ares := <-accCh
	require.NoError(t, ares.err)
	require.NotEmpty(t, ares.cn)

	err = ares.cn.Close()
	require.NoError(t, err)
}
