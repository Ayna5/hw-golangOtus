package main

import (
	"flag"
	"fmt"
	"net"
	"os"
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

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Errorf("cannot connect: %w", err)
	}
	defer client.Close()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Errorf("can't receieve: %v", err)
			return
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Errorf("can't send: %v", err)
			return
		}
	}()
}
