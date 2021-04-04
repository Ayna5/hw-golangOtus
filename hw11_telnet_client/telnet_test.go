package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client, err := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out, func() {})
			require.NoError(t, err)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
	t.Run("in is nil", func(t *testing.T) {
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client, err := NewTelnetClient("127.0.0.1:9000", timeout, nil, out, func() {})
		require.NoError(t, err)
		require.Equal(t, "in is nil", client.Connect().Error())
	})

	t.Run("out is nil", func(t *testing.T) {
		in := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client, err := NewTelnetClient("127.0.0.1:9000", timeout, ioutil.NopCloser(in), nil, func() {})
		require.NoError(t, err)
		require.Equal(t, "out is nil", client.Connect().Error())
	})

	t.Run("EOF", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			pr, pw := io.Pipe()
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client, err := NewTelnetClient(l.Addr().String(), timeout, pr, out, func() {})
			require.NoError(t, err)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			err = pw.Close()
			require.NoError(t, err)

			_, err = pr.Read(in.Bytes())
			require.Contains(t, err.Error(), "EOF")

			err = client.Send()
			require.NoError(t, err) // тут все равно не приходит ошибка, либо я что-то не то делаю
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
		}()

		wg.Wait()
	})
}
