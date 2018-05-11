package mocktrans

import (
	"io"
	"net"
)

type transport struct {
	acceptQueue chan net.Conn
	addr        net.Addr
}

type transportOpt func(l *transport) error

func NewTransport(opts ...transportOpt) (tr *transport, err error) {
	tr = &transport{
		acceptQueue: make(chan net.Conn, 5),
		addr:        fakeAddr{},
	}
	for i := range opts {
		err = opts[i](tr)
		if err != nil {
			tr = nil
			return
		}
	}
	if tr.addr == nil {
		tr.addr = newFakeAddr()
	}
	return
}

func (tr *transport) Dial() (net.Conn, error) {
	srvRead, cliWrite := io.Pipe()
	cliRead, srvWrite := io.Pipe()

	srv, err := NewConn(ConnWriter(srvWrite), ConnReader(srvRead))
	if err != nil {
		return nil, err
	}

	cli, err := NewConn(ConnWriter(cliWrite), ConnReader(cliRead))
	if err != nil {
		return nil, err
	}
	tr.acceptQueue <- srv
	return cli, nil
}

func (tr *transport) Accept() (net.Conn, error) {
	return <-tr.acceptQueue, nil
}

func (*transport) Close() error {
	return nil
}

func (tr *transport) Addr() net.Addr {
	return tr.addr
}
