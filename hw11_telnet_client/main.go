package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout, default:10s")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		return
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	addr := net.JoinHostPort(host, port)

	telnetClient := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	err := telnetClient.Connect()
	if err != nil {
		fmt.Println("Connect error: ", err)
		return
	}
	defer telnetClient.Close()

	errChan := make(chan error, 2)

	go func() {
		errChan <- telnetClient.Receive()
	}()

	go func() {
		errChan <- telnetClient.Send()
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "\nInterrupted by User")
	case err := <-errChan:
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr, "\nConnection closed")
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "\nI/O error: %v\n", err)
		}
	}
}
