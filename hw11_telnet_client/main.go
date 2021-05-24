package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	if len(flag.Args())<2 {
		log.Fatal("host and port must be provided as two arguments")
	}
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client, err := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout, cancel)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot execute NewTelnetClient: %w", err))
		return
	}

	if err := client.Connect(); err != nil {
		fmt.Println(fmt.Errorf("cannot connect: %w", err))
		return
	}
	defer client.Close()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Println(fmt.Errorf("can't receieve: %w", err))
			return
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Println(fmt.Errorf("can't send: %w", err))
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
