package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/praveenmahasena/server/internal/listner"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cancelCh := make(chan os.Signal, 1)

		signal.Notify(cancelCh, os.Interrupt, syscall.SIGKILL, syscall.SIGINT) // the same thing

		select {
		case sig := <-cancelCh:
			fmt.Printf("Received signal: %s\n", sig)
			cancel()
		case <-ctx.Done():
			fmt.Printf("closing signal goroutine\n")
		}
	}()

	l := listner.New(ctx, "tcp", ":42069")
	return l.Run()
}
