package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	timeout    time.Duration
	host, port string
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	flag.StringVar(&host, "host", "localhost", "host")
	flag.StringVar(&port, "port", "9000", "port")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout, cancel)

	if err := client.Connect(); err != nil {
		fmt.Println(fmt.Errorf("cannot connect: %w", err))
	}
	defer client.Close()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Errorf("can't receieve: %w", err)
			return
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Errorf("can't send: %w", err)
			return
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
		close(sig)
	}
}
