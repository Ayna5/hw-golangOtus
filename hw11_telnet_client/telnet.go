package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	cancel  context.CancelFunc
}

func (t *telnetClient) Connect() error {
	if t.in == nil {
		return errors.New("in is nil")
	}
	if t.out == nil {
		return errors.New("out is nil")
	}
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}
	t.conn = conn
	return nil
}

func (t *telnetClient) Send() error {
	defer t.cancel()
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return fmt.Errorf("cannot send: %w", err)
	}
	return nil
}

func (t *telnetClient) Receive() error {
	defer t.cancel()
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return fmt.Errorf("cannot receive: %w", err)
	}
	return nil
}

func (t *telnetClient) Close() error {
	return t.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer, cancel context.CancelFunc) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		cancel:  cancel,
	}
}
