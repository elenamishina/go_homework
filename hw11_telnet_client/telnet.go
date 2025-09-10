package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	addr    string
	timeout time.Duration
	conn    net.Conn
	in      *bufio.Reader
	out     *bufio.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{addr: address, timeout: timeout, in: bufio.NewReader(in), out: bufio.NewWriter(out)}
}

func (t *Telnet) Connect() error {
	dialer := &net.Dialer{
		Timeout: t.timeout,
	}
	conn, err := dialer.DialContext(context.Background(), "tcp", t.addr)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if tc, ok := t.conn.(*net.TCPConn); ok {
		_ = tc.CloseWrite()
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (t *Telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (t *Telnet) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		return err
	}
	return nil
}
