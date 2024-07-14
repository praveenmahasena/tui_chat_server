package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/praveenmahasena/server/internal/listener"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cancelCh := make(chan os.Signal, 1)

		signal.Notify(cancelCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // the same thing

		select {
		case sig := <-cancelCh:
			fmt.Printf("Received signal: %s\n", sig)
			cancel()
		case <-ctx.Done():
			fmt.Printf("closing signal goroutine\n")
		}
	}()

	l := listener.New(ctx, "tcp", ":42069")
	return l.Run()
}
