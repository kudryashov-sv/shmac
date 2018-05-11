package mocktrans

import (
	"net"
	"time"

	"sync"

	"io"

	"github.com/pkg/errors"
)

type mockConn struct {
	mtx sync.Mutex

	readDeadline  time.Time
	writeDeadline time.Time

	localAddr  net.Addr
	remoteAddr net.Addr

	reader io.ReadCloser
	writer io.WriteCloser
}

type mockConnOpt func(cn *mockConn) error

func ConnReader(r io.ReadCloser) mockConnOpt {
	return func(cn *mockConn) error {
		cn.reader = r
		return nil
	}
}

func ConnWriter(w io.WriteCloser) mockConnOpt {
	return func(cn *mockConn) error {
		cn.writer = w
		return nil
	}
}

func NewConn(opts ...mockConnOpt) (cn *mockConn, err error) {
	cn = new(mockConn)
	for i := range opts {
		err = opts[i](cn)
		if err != nil {
			cn = nil
			return
		}
	}
	cn.localAddr = cn.maybeAddr(cn.localAddr)
	cn.remoteAddr = cn.maybeAddr(cn.remoteAddr)

	if cn.reader == nil && cn.writer == nil {
		cn.reader, cn.writer = io.Pipe()
	} else if (cn.reader == nil && cn.writer != nil) || (cn.reader != nil && cn.writer == nil) {
		err = errors.Errorf("internal error, reader: %#v, writer: %#v", cn.reader, cn.writer)
		return
	}
	return
}

func (mc *mockConn) Read(b []byte) (n int, err error) {
	mc.mtx.Lock()

	if !mc.readDeadline.IsZero() && mc.readDeadline.Before(time.Now()) {
		mc.mtx.Unlock()
		err = NewError("read timeout", true, true)
		return
	}
	mc.mtx.Unlock()

	return mc.reader.Read(b)
}

func (mc *mockConn) Write(b []byte) (n int, err error) {
	mc.mtx.Lock()
	if !mc.writeDeadline.IsZero() && mc.writeDeadline.Before(time.Now()) {
		mc.mtx.Unlock()
		err = NewError("write timeout", true, true)
		return
	}
	mc.mtx.Unlock()

	return mc.writer.Write(b)
}

func (mc *mockConn) Close() error {
	re := mc.reader.Close()
	we := mc.writer.Close()
	if re != nil || we != nil {
		return errors.Errorf("close reader: %v, close writer: %v", re, we)
	}
	return nil
}

func (mc *mockConn) LocalAddr() net.Addr {
	return mc.localAddr
}

func (mc *mockConn) RemoteAddr() net.Addr {
	return mc.remoteAddr
}

func (mc *mockConn) SetDeadline(t time.Time) error {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()

	mc.readDeadline = t
	mc.writeDeadline = t

	return nil
}

func (*mockConn) SetReadDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (*mockConn) SetWriteDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (mc *mockConn) maybeAddr(addr net.Addr) net.Addr {
	if addr == nil {
		return newFakeAddr()
	}
	return addr
}
