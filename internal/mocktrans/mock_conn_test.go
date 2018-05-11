package mocktrans

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewConn(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)
	require.Implements(t, (*net.Conn)(nil), cn)
}

func TestMockConn_ReadWrite(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)

	exp := []byte("test message")
	act := make([]byte, len(exp))
	writeRes := make(chan ioRes)
	readRes := make(chan ioRes)
	go func() {
		res := ioRes{}
		res.n, res.err = cn.Write(exp)
		writeRes <- res
	}()

	go func() {
		res := ioRes{}
		res.n, res.err = cn.Read(act)
		readRes <- res
	}()

	wr := waitResult(t, writeRes, "wait write timeout")
	rr := waitResult(t, readRes, "wait read timeout")

	require.NoError(t, wr.err)
	require.Equal(t, len(exp), wr.n)

	require.NoError(t, rr.err)
	require.Equal(t, len(exp), rr.n)
	require.Equal(t, exp, act)
}

func TestMockConn_SetDeadline(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)

	err = cn.SetDeadline(time.Now())
	require.NoError(t, err)

	n, err := cn.Read(nil)
	require.Error(t, err)
	require.Zero(t, n)
	nerr, ok := err.(net.Error)
	require.True(t, ok)
	require.True(t, nerr.Timeout())

	n, err = cn.Write(nil)
	require.Error(t, err)
	require.Zero(t, n)
	nerr, ok = err.(net.Error)
	require.True(t, ok)
	require.True(t, nerr.Timeout())
}

func TestMockConn_RemoteAddr(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)

	addr := cn.RemoteAddr()
	require.Equal(t, newFakeAddr(), addr)
}

func TestMockConn_LocalAddr(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)

	addr := cn.LocalAddr()
	require.Equal(t, newFakeAddr(), addr)
}

func TestMockConn_Close(t *testing.T) {
	cn, err := NewConn()
	require.NoError(t, err)

	err = cn.Close()
	require.NoError(t, err)
}

type ioRes struct {
	n   int
	err error
}

func waitResult(t *testing.T, ch <-chan ioRes, format string, args ...interface{}) (res ioRes) {
	select {
	case res = <-ch:
		return
	case <-time.After(time.Second * 3):
		t.Fatalf(format, args...)
	}
	return
}
