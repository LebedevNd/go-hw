package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const defaultTimeout = 10

func init() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout*time.Second, "connection timeout")
}

func main() {
	host := os.Args[2]
	port := os.Args[3]
	address := host + ":" + port
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		defer cancelFunc()
		err := client.Send()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		defer cancelFunc()
		err := client.Receive()
		if err != nil {
			fmt.Println(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		cancelFunc()
		signal.Stop(sigCh)
		return

	case <-ctx.Done():
		close(sigCh)
		return
	}
}
